//go:build integration

package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/etc"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/job"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/persistence/redis"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/redisx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestStore is an integration test for the Redis persistence store.
func TestStore(t *testing.T) {
	if testing.Short() {
		t.Skip("An integration test")
	}

	ctx := context.Background()
	redisC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "redis:5.0.5",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err, "should start redis container")
	defer func() {
		_ = redisC.Terminate(ctx)
	}()

	redisURL := getRedisURL(t, ctx, redisC)

	pool, err := redisx.NewPool(etc.RedisPool{
		URL: redisURL,
	})
	require.NoError(t, err)

	store := redis.NewStore(etc.RedisStore{
		Namespace:  "harbor.scanner.khulnasoft:store",
		ScanJobTTL: parseDuration(t, "10s"),
	}, pool)

	t.Run("CRUD", func(t *testing.T) {
		scanJobID := "123"

		err := store.Create(job.ScanJob{
			ID:     scanJobID,
			Status: job.Pending,
		})
		require.NoError(t, err, "saving scan job should not fail")

		j, err := store.Get(scanJobID)
		require.NoError(t, err, "getting scan job should not fail")
		assert.Equal(t, &job.ScanJob{
			ID:     scanJobID,
			Status: job.Pending,
		}, j)

		err = store.UpdateStatus(scanJobID, job.Running)
		require.NoError(t, err, "updating scan job status should not fail")

		j, err = store.Get(scanJobID)
		require.NoError(t, err, "getting scan job should not fail")
		assert.Equal(t, &job.ScanJob{
			ID:     scanJobID,
			Status: job.Running,
		}, j)

		scanReport := harbor.ScanReport{
			Severity: harbor.SevHigh,
			Vulnerabilities: []harbor.VulnerabilityItem{
				{
					ID:         "CVE-2013-1400",
					Pkg:        "openssl",
					Version:    "2.4",
					FixVersion: "2.4.2",
					Severity:   harbor.SevHigh,
					Links: []string{
						"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2013-1400",
					},
				},
			},
		}

		err = store.UpdateReport(scanJobID, scanReport)
		require.NoError(t, err, "updating scan job reports should not fail")

		j, err = store.Get(scanJobID)
		require.NoError(t, err, "retrieving scan job should not fail")
		require.NotNil(t, j, "retrieved scan job must not be nil")
		assert.Equal(t, scanReport, j.Report)

		err = store.UpdateStatus(scanJobID, job.Finished)
		require.NoError(t, err)

		time.Sleep(parseDuration(t, "12s"))

		j, err = store.Get(scanJobID)
		require.NoError(t, err, "retrieve scan job should not fail")
		require.Nil(t, j, "retrieved scan job should be nil, i.e. expired")
	})
}

func getRedisURL(t *testing.T, ctx context.Context, redisC tc.Container) string {
	t.Helper()
	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)
	return fmt.Sprintf("redis://%s:%d", host, port.Int())
}

func parseDuration(t *testing.T, s string) time.Duration {
	t.Helper()
	d, err := time.ParseDuration(s)
	require.NoError(t, err, "should parse duration %s", s)
	return d
}
