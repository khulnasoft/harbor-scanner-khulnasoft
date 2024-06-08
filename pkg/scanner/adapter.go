package scanner

import (
	"fmt"

	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
)

type Adapter interface {
	Scan(req harbor.ScanRequest) (harbor.ScanReport, error)
}

type adapter struct {
	command     khulnasoft.Command
	transformer Transformer
}

func NewAdapter(command khulnasoft.Command, transformer Transformer) Adapter {
	return &adapter{
		command:     command,
		transformer: transformer,
	}
}

func (s *adapter) Scan(req harbor.ScanRequest) (harborReport harbor.ScanReport, err error) {
	username, password, err := req.Registry.GetBasicCredentials()
	if err != nil {
		err = fmt.Errorf("getting basic credentials from scan request: %w", err)
		return
	}

	khulnasoftReport, err := s.command.Scan(khulnasoft.ImageRef{
		Repository: req.Artifact.Repository,
		Tag:        req.Artifact.Tag,
		Digest:     req.Artifact.Digest,
		Auth: khulnasoft.RegistryAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		return
	}
	harborReport = s.transformer.Transform(req.Artifact, khulnasoftReport)
	return
}
