package scanner

import (
	"bytes"
	"encoding/json"

	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/models"
)

type vulnDetails struct {
	Pairs []vulnerabilityPackagePair `json:"details"`
}

type vulnerabilityPackagePair struct {
	Vulnerability models.Vulnerability `json:"vulnerability"`
	Package       pkg.Package          `json:"package"`
}

func PrintVulnDetails(cfg models.PresenterConfig) (string, error) {
	vulnDetails := vulnDetails{}
	for _, match := range cfg.Matches.Sorted() {
		vuln := models.NewVulnerability(match.Vulnerability, nil)
		pair := vulnerabilityPackagePair{
			Vulnerability: vuln,
			Package:       match.Package,
		}
		vulnDetails.Pairs = append(vulnDetails.Pairs, pair)
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
