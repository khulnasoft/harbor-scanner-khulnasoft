package persistence

import (
	"context"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/job"
)

type Store interface {
	Create(ctx context.Context, scanJob job.ScanJob) error
	Get(ctx context.Context, scanJobKey job.ScanJobKey) (*job.ScanJob, error)
	UpdateStatus(ctx context.Context, scanJobKey job.ScanJobKey, newStatus job.ScanJobStatus, error ...string) error
	UpdateReport(ctx context.Context, scanJobKey job.ScanJobKey, report harbor.ScanReport) error
}
