package mock

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/http/api"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/trivy"
	"github.com/stretchr/testify/mock"
)

type Transformer struct {
	mock.Mock
}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) Transform(mediaType api.MediaType, req harbor.ScanRequest, source trivy.Report) harbor.ScanReport {
	args := t.Called(mediaType, req, source)
	return args.Get(0).(harbor.ScanReport)
}
