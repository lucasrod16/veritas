package scanner

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
	"github.com/anchore/grype/grype/matcher"
	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/cyclonedx"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/anchore/grype/grype/store"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/sbom"
	"golang.org/x/sync/errgroup"
)

const grypeDBListingURL = "https://toolbox-data.anchore.io/grype/databases/listing.json"

var grypeDBdir = filepath.Join(xdg.CacheHome, "veritas", "db")

func Scan(userInput string) (string, error) {
	var (
		err error
		g   errgroup.Group

		store  *store.Store
		status *db.Status
		closer *db.Closer

		packages   []pkg.Package
		pkgContext pkg.Context
		sbom       *sbom.SBOM
	)

	g.Go(func() error {
		store, status, closer, err = grype.LoadVulnerabilityDB(db.Config{
			DBRootDir:  grypeDBdir,
			ListingURL: grypeDBListingURL,
		}, true)
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		packages, pkgContext, sbom, err = pkg.Provide(userInput, pkg.ProviderConfig{
			SyftProviderConfig: pkg.SyftProviderConfig{
				SBOMOptions: syft.DefaultCreateSBOMConfig(),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to catalog: %w", err)
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		return "", err
	}

	defer closer.Close()

	vulnMatcher := grype.VulnerabilityMatcher{
		Store:    *store,
		Matchers: matcher.NewDefaultMatchers(matcher.Config{}),
	}

	matches, _, err := vulnMatcher.FindMatches(packages, pkgContext)
	if err != nil {
		return "", err
	}

	pres := cyclonedx.NewJSONPresenter(models.PresenterConfig{
		Matches:          *matches,
		Packages:         packages,
		Context:          pkgContext,
		MetadataProvider: store,
		SBOM:             sbom,
		DBStatus:         status,
	})

	buf := &bytes.Buffer{}
	err = pres.Present(buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
