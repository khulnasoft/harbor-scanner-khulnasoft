package scanner

import (
	"testing"

	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var NoError error

func TestAdapter_Scan(t *testing.T) {

	t.Run("Should return error when getting registry credentials fails", func(t *testing.T) {
		command := &khulnasoft.MockCommand{}
		transformer := &MockTransformer{}

		scanRequest := harbor.ScanRequest{
			Registry: harbor.Registry{
				Authorization: "Bearer 0123456789",
			},
		}

		_, err := NewAdapter(command, transformer).Scan(scanRequest)
		assert.EqualError(t, err, "getting basic credentials from scan request: unsupported authorization type: Bearer")

		command.AssertExpectations(t)
		transformer.AssertExpectations(t)
	})

	t.Run("Should return Harbor report", func(t *testing.T) {
		command := &khulnasoft.MockCommand{}
		transformer := &MockTransformer{}

		artifact := harbor.Artifact{
			Repository: "library/golang",
			Tag:        "1.12.4",
		}
		scanRequest := harbor.ScanRequest{
			Registry: harbor.Registry{
				Authorization: "Basic cm9ib3ROYW1lOnJvYm90UGFzc3dvcmQ=",
			},
			Artifact: artifact,
		}
		imageRef := khulnasoft.ImageRef{
			Repository: "library/golang",
			Tag:        "1.12.4",
			Auth: khulnasoft.RegistryAuth{
				Username: "robotName",
				Password: "robotPassword",
			},
		}

		khulnasoftReport := khulnasoft.ScanReport{}
		harborReport := harbor.ScanReport{}

		command.On("Scan", imageRef).Return(khulnasoftReport, NoError)
		transformer.On("Transform", artifact, khulnasoftReport).Return(harborReport)

		r, err := NewAdapter(command, transformer).Scan(scanRequest)
		require.NoError(t, err)
		require.Equal(t, harborReport, r)

		command.AssertExpectations(t)
		transformer.AssertExpectations(t)
	})

}
