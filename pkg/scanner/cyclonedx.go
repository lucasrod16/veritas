package scanner

import (
	"bytes"

	"github.com/anchore/grype/grype/presenter/cyclonedx"
	"github.com/anchore/grype/grype/presenter/models"
)

func PrintCycloneDXJSON(cfg models.PresenterConfig) (string, error) {
	pres := cyclonedx.NewJSONPresenter(cfg)
	buf := &bytes.Buffer{}
	err := pres.Present(buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
