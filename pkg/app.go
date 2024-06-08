package pkg

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/etc"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/ext"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/http/api"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/http/api/v1"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/persistence/redis"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/redisx"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/scanner"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/work"
	log "github.com/sirupsen/logrus"
)

func Run(info etc.BuildInfo) error {
	log.WithFields(log.Fields{
		"version":  info.Version,
		"commit":   info.Commit,
		"built_at": info.Date,
	}).Info("Starting harbor-scanner-khulnasoft")

	config, err := etc.GetConfig()
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	if _, err := os.Stat(config.KhulnasoftCSP.ReportsDir); os.IsNotExist(err) {
		log.WithField("path", config.KhulnasoftCSP.ReportsDir).Debug("Creating reports dir")
		err = os.MkdirAll(config.KhulnasoftCSP.ReportsDir, os.ModeDir)
		if err != nil {
			return fmt.Errorf("creating reports dir: %w", err)
		}
	}

	pool, err := redisx.NewPool(config.RedisPool)
	if err != nil {
		return fmt.Errorf("constructing connection pool: %w", err)
	}

	workPool := work.New()
	command := khulnasoft.NewCommand(config.KhulnasoftCSP, ext.DefaultAmbassador)
	transformer := scanner.NewTransformer(ext.NewSystemClock())
	adapter := scanner.NewAdapter(command, transformer)
	store := redis.NewStore(config.RedisStore, pool)
	enqueuer := scanner.NewEnqueuer(workPool, adapter, store)
	apiServer := api.NewServer(config.API, v1.NewAPIHandler(info, config, enqueuer, store))

	shutdownComplete := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		captured := <-sigint
		log.WithField("signal", captured.String()).Debug("Trapped os signal")

		apiServer.Shutdown()
		workPool.Shutdown()

		close(shutdownComplete)
	}()

	workPool.Start()
	apiServer.ListenAndServe()

	<-shutdownComplete
	return nil
}
