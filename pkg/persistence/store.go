package persistence

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/job"
)

// Store defines methods for persisting ScanJobs and associated ScanReports.
type Store interface {
	Create(scanJob job.ScanJob) error
	Get(scanJobID string) (*job.ScanJob, error)
	UpdateStatus(scanJobID string, newStatus job.Status, error ...string) error
	UpdateReport(scanJobID string, reports harbor.ScanReport) error
}
