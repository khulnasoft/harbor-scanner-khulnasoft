package scanner

import (
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/khulnasoft"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/etc"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/ext"
	"github.com/khulnasoft/harbor-scanner-khulnasoft/pkg/harbor"
	log "github.com/sirupsen/logrus"
)

type Transformer interface {
	Transform(artifact harbor.Artifact, source khulnasoft.ScanReport) harbor.ScanReport
}

func NewTransformer(clock ext.Clock) Transformer {
	return &transformer{
		clock: clock,
	}
}

type transformer struct {
	clock ext.Clock
}

func (t *transformer) Transform(artifact harbor.Artifact, source khulnasoft.ScanReport) harbor.ScanReport {
	log.WithFields(log.Fields{
		"digest":          source.Digest,
		"image":           source.Image,
		"summary":         source.Summary,
		"scan_options":    source.ScanOptions,
		"changed_results": source.ChangedResults,
		"partial_results": source.PartialResults,
	}).Debug("Transforming scan report")
	var items []harbor.VulnerabilityItem

	for _, resourceScan := range source.Resources {
		var pkg string
		switch resourceScan.Resource.Type {
		case khulnasoft.Library:
			pkg = resourceScan.Resource.Path
		case khulnasoft.Package:
			pkg = resourceScan.Resource.Name
		default:
			pkg = resourceScan.Resource.Name
		}

		for _, vln := range resourceScan.Vulnerabilities {
			items = append(items, harbor.VulnerabilityItem{
				ID:          vln.Name,
				Pkg:         pkg,
				Version:     resourceScan.Resource.Version,
				FixVersion:  vln.FixVersion,
				Severity:    t.getHarborSeverity(vln),
				Description: vln.Description,
				Links:       t.toLinks(vln),
			})
		}
	}

	return harbor.ScanReport{
		GeneratedAt:     t.clock.Now(),
		Scanner:         etc.GetScannerMetadata(),
		Artifact:        artifact,
		Severity:        t.getHighestSeverity(items),
		Vulnerabilities: items,
	}
}

func (t *transformer) getHarborSeverity(v khulnasoft.Vulnerability) harbor.Severity {
	var severity harbor.Severity
	switch v.KhulnasoftSeverity {
	case "critical":
		severity = harbor.SevCritical
	case "high":
		severity = harbor.SevHigh
	case "medium":
		severity = harbor.SevMedium
	case "low":
		severity = harbor.SevLow
	case "negligible":
		severity = harbor.SevNegligible
	default:
		log.WithField("severity", v.KhulnasoftSeverity).Warn("Unknown Khulnasoft severity")
		severity = harbor.SevUnknown
	}
	return severity
}

func (t *transformer) toLinks(v khulnasoft.Vulnerability) []string {
	var links []string
	if v.NVDURL != "" {
		links = append(links, v.NVDURL)
	}
	if v.VendorURL != "" {
		links = append(links, v.VendorURL)
	}
	return links
}

func (t *transformer) getHighestSeverity(items []harbor.VulnerabilityItem) (highest harbor.Severity) {
	highest = harbor.SevUnknown

	for _, v := range items {
		if v.Severity > highest {
			highest = v.Severity
		}
	}

	return
}
