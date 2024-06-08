package scanner

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/stretchr/testify/mock"
)

type MockEnqueuer struct {
	mock.Mock
}

func (e *MockEnqueuer) Enqueue(request harbor.ScanRequest) (string, error) {
	args := e.Called(request)
	return args.String(0), args.Error(1)
}
