package printer

import (
	"bytes"
	"encoding/json"

	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/anchore/grype/grype/vulnerability"
)

type vulnDetail struct {
	Vulnerability vulnerability.Vulnerability `json:"vulnerability"`
	Package       pkg.Package                 `json:"package"`
	Severity      string                      `json:"severity"`
}

func NewVulnDetailsPrinter() vulnDetail { return vulnDetail{} }

func (vd vulnDetail) Print(cfg models.PresenterConfig) (string, error) {
	vulnDetails := []vulnDetail{}
	for _, match := range cfg.Matches.Sorted() {
		vulnMetadata, err := cfg.MetadataProvider.GetMetadata(match.Vulnerability.ID, match.Vulnerability.Namespace)
		if err != nil {
			return "", err
		}
		vd.Vulnerability = match.Vulnerability
		vd.Package = match.Package
		vd.Severity = vulnMetadata.Severity
		vulnDetails = append(vulnDetails, vd)
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

var _ Printer = (*vulnDetail)(nil)
