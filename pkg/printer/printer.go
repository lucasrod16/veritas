package printer

import "github.com/anchore/grype/grype/presenter/models"

type Printer interface {
	Print(cfg models.PresenterConfig) (string, error)
}
