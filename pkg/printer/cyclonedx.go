package printer

import (
	"bytes"

	"github.com/anchore/grype/grype/presenter/cyclonedx"
	"github.com/anchore/grype/grype/presenter/models"
)

type cycloneDX struct{}

func NewCycloneDXPrinter() cycloneDX { return cycloneDX{} }

func (cycloneDX) Print(cfg models.PresenterConfig) (string, error) {
	pres := cyclonedx.NewJSONPresenter(cfg)
	buf := &bytes.Buffer{}
	err := pres.Present(buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var _ Printer = (*cycloneDX)(nil)
