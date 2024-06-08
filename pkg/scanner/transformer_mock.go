package scanner

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/stretchr/testify/mock"
)

type MockTransformer struct {
	mock.Mock
}

func (t *MockTransformer) Transform(artifact harbor.Artifact, source khulnasoft.ScanReport) harbor.ScanReport {
	args := t.Called(artifact, source)
	return args.Get(0).(harbor.ScanReport)
}
