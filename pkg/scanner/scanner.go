package scanner

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/anchore/grype/grype"
	"github.com/anchore/grype/grype/db"
	"github.com/anchore/grype/grype/matcher"
	"github.com/anchore/grype/grype/matcher/golang"
	"github.com/anchore/grype/grype/matcher/java"
	"github.com/anchore/grype/grype/matcher/stock"
	"github.com/anchore/grype/grype/pkg"
	"github.com/anchore/grype/grype/presenter/cyclonedx"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/anchore/grype/grype/store"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/sbom"
	"golang.org/x/sync/errgroup"
)

const mavenSearchBaseURL = "https://search.maven.org/solrsearch/select"

func Scan(userInput string) (string, error) {
	var g errgroup.Group
	var err error

	var store *store.Store
	var status *db.Status
	var closer *db.Closer
	g.Go(func() error {
		store, status, closer, err = grype.LoadVulnerabilityDB(newGrypeDBCfg(), true)
		if err != nil {
			return err
		}
		return nil
	})

	var packages []pkg.Package
	var pkgContext pkg.Context
	var sbom *sbom.SBOM
	g.Go(func() error {
		packages, pkgContext, sbom, err = pkg.Provide(userInput, getProviderConfig())
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
		Matchers: getMatchers(),
	}

	remainingMatches, _, err := vulnMatcher.FindMatches(packages, pkgContext)
	if err != nil {
		return "", err
	}

	presenterCfg := models.PresenterConfig{
		Matches:          *remainingMatches,
		Packages:         packages,
		Context:          pkgContext,
		MetadataProvider: store,
		SBOM:             sbom,
		DBStatus:         status,
	}

	pres := cyclonedx.NewJSONPresenter(presenterCfg)
	buf := &bytes.Buffer{}
	err = pres.Present(buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newGrypeDBCfg() db.Config {
	return db.Config{
		DBRootDir:  filepath.Join(xdg.CacheHome, "veritas", "db"),
		ListingURL: "https://toolbox-data.anchore.io/grype/databases/listing.json",
	}
}

func getProviderConfig() pkg.ProviderConfig {
	return pkg.ProviderConfig{
		SyftProviderConfig: pkg.SyftProviderConfig{
			SBOMOptions: syft.DefaultCreateSBOMConfig(),
		},
	}
}

func getMatchers() []matcher.Matcher {
	return matcher.NewDefaultMatchers(
		matcher.Config{
			Golang: golang.MatcherConfig{
				AlwaysUseCPEForStdlib:                  true,
				AllowMainModulePseudoVersionComparison: false,
			},
			Java: java.MatcherConfig{
				ExternalSearchConfig: java.ExternalSearchConfig{
					SearchMavenUpstream: true,
					MavenBaseURL:        mavenSearchBaseURL,
				},
			},
			Stock: stock.MatcherConfig{UseCPEs: true},
		},
	)
}
