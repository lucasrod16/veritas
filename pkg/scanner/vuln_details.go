package scanner

import (
	"bytes"
	"encoding/json"

	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/anchore/grype/grype/vulnerability"
)

type vulnData struct {
	Vulnerability vulnerability.Vulnerability `json:"vulnerability"`
	Package       pkg.Package                 `json:"package"`
	Severity      string                      `json:"severity"`
}

func PrintVulnDetails(cfg models.PresenterConfig) (string, error) {
	vulnDetails := []vulnData{}
	for _, match := range cfg.Matches.Sorted() {
		vulnMetadata, err := cfg.MetadataProvider.GetMetadata(match.Vulnerability.ID, match.Vulnerability.Namespace)
		if err != nil {
			return "", err
		}
		vulnDetails = append(vulnDetails, vulnData{
			Vulnerability: match.Vulnerability,
			Package:       match.Package,
			Severity:      vulnMetadata.Severity,
		})
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(vulnDetails)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
