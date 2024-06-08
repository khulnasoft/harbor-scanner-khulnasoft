package scanner

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/ext"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransformer_Transform(t *testing.T) {
	now := time.Now()

	artifact := harbor.Artifact{
		Repository: "library/golang",
		Tag:        "1.12.4",
	}

	khulnasoftReport := khulnasoft.ScanReport{
		Resources: []khulnasoft.ResourceScan{
			{
				Resource: khulnasoft.Resource{
					Type:    khulnasoft.Package,
					Name:    "openssl",
					Version: "2.8.3",
				},
				Vulnerabilities: []khulnasoft.Vulnerability{
					{
						Name:         "CVE-0001-0020",
						KhulnasoftSeverity: "high",
						NVDURL:       "http://nvd?id=CVE-0001-0020",
					},
					{
						Name:         "CVE-3045-2011",
						KhulnasoftSeverity: "low",
					},
				},
			},
			{
				Resource: khulnasoft.Resource{
					Type: khulnasoft.Library,
					Path: "/app/main.rb",
				},
				Vulnerabilities: []khulnasoft.Vulnerability{
					{
						Name:         "CVE-9900-1100",
						KhulnasoftSeverity: "critical",
					},
				},
			},
		},
	}

	harborReport := NewTransformer(ext.NewFixedClock(now)).Transform(artifact, khulnasoftReport)
	assert.Equal(t, harbor.ScanReport{
		GeneratedAt: now,
		Artifact: harbor.Artifact{
			Repository: "library/golang",
			Tag:        "1.12.4",
		},
		Scanner: harbor.Scanner{
			Name:    "Khulnasoft Enterprise",
			Vendor:  "Khulnasoft Security",
			Version: "Unknown",
		},
		Severity: harbor.SevCritical,
		Vulnerabilities: []harbor.VulnerabilityItem{
			{
				ID:       "CVE-0001-0020",
				Pkg:      "openssl",
				Version:  "2.8.3",
				Severity: harbor.SevHigh,
				Links: []string{
					"http://nvd?id=CVE-0001-0020",
				},
			},
			{
				ID:       "CVE-3045-2011",
				Pkg:      "openssl",
				Version:  "2.8.3",
				Severity: harbor.SevLow,
				Links:    ([]string)(nil),
			},
			{
				ID:       "CVE-9900-1100",
				Pkg:      "/app/main.rb",
				Version:  "",
				Severity: harbor.SevCritical,
				Links:    ([]string)(nil),
			},
		},
	}, harborReport)
}
