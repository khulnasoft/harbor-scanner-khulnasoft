package main

import (
	"os"

	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/etc"
	log "github.com/sirupsen/logrus"
)

var (
	// Default wise GoReleaser sets three ldflags:
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(etc.GetLogLevel())
	log.SetReportCaller(false)
	log.SetFormatter(&log.JSONFormatter{})

	info := etc.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}

	if err := pkg.Run(info); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
