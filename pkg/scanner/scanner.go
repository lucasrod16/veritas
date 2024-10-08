package scanner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
	"github.com/anchore/grype/grype/db/legacy/distribution"
	"github.com/anchore/grype/grype/matcher"
	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/anchore/grype/grype/store"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/sbom"
	"golang.org/x/sync/errgroup"
)

const grypeDBListingURL = "https://toolbox-data.anchore.io/grype/databases/listing.json"

var grypeDBdir = filepath.Join(xdg.CacheHome, "veritas", "db")

func Scan(userInput string) (models.PresenterConfig, *db.Closer, error) {
	var err error
	var g errgroup.Group

	var store *store.Store
	var status *distribution.Status
	var closer *db.Closer

	g.Go(func() error {
		store, status, closer, err = grype.LoadVulnerabilityDB(distribution.Config{
			DBRootDir:  grypeDBdir,
			ListingURL: grypeDBListingURL,
		}, true)
		if err != nil {
			return err
		}
		return nil
	})

	var packages []pkg.Package
	var pkgContext pkg.Context
	var sbom *sbom.SBOM

	g.Go(func() error {
		packages, pkgContext, sbom, err = pkg.Provide(userInput, pkg.ProviderConfig{
			SyftProviderConfig: pkg.SyftProviderConfig{
				SBOMOptions:            syft.DefaultCreateSBOMConfig(),
				DefaultImagePullSource: "registry",
			},
		})

		if err == nil {
			return nil
		}

		switch {
		case strings.Contains(err.Error(), "failed to get image descriptor from registry"):
			return fmt.Errorf(
				"unable to scan %q. please ensure the image exists, is publicly accessible, and does not require authentication",
				userInput,
			)
		case strings.Contains(err.Error(), "unable to parse registry reference"):
			return fmt.Errorf("unable to scan %q. please ensure you provide a valid image reference", userInput)
		default:
			return err
		}
	})

	err = g.Wait()
	if err != nil {
		return models.PresenterConfig{}, nil, err
	}

	vulnMatcher := grype.VulnerabilityMatcher{
		Store:    *store,
		Matchers: matcher.NewDefaultMatchers(matcher.Config{}),
	}

	matches, _, err := vulnMatcher.FindMatches(packages, pkgContext)
	if err != nil {
		return models.PresenterConfig{}, nil, err
	}

	return models.PresenterConfig{
			Matches:          *matches,
			Packages:         packages,
			Context:          pkgContext,
			MetadataProvider: store,
			SBOM:             sbom,
			DBStatus:         status,
		},
		closer,
		nil
}
