package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

const (
	bomFile    = "bom.json"
	schemaFile = "bom-1.6.schema.json"

	scanURL = "http://localhost:8080/scan?image="

	// images without labels result in non-compliant SBOMs currently.
	// This PR fixes it: https://github.com/anchore/syft/pull/3119.
	// Waiting on syft to release the fix.
	image = "ubuntu:latest"

	payloadEndpoint = scanURL + image
)

var projectPath = filepath.Join("cmd", "schema-validator")

func main() {
	cp, err := currentPath()
	if err != nil {
		log.Fatal(err)
	}

	err = generateSBOM(cp)
	if err != nil {
		log.Fatal(err)
	}

	bom, err := os.ReadFile(filepath.Join(cp, bomFile))
	if err != nil {
		log.Fatal(err)
	}
	schema, err := os.ReadFile(filepath.Join(cp, schemaFile))
	if err != nil {
		log.Fatal(err)
	}

	schemaLoader := gojsonschema.NewStringLoader(string(schema))
	documentLoader := gojsonschema.NewStringLoader(string(bom))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatal(err)
	}

	if result.Valid() {
		fmt.Println("✅ CycloneDX JSON SBOM is valid ✅")
	} else {
		fmt.Println("❌ CycloneDX JSON SBOM is NOT valid ❌")
		for _, schemaErr := range result.Errors() {
			fmt.Printf("- %s\n", schemaErr)
		}
		os.Exit(1)
	}
}

func generateSBOM(absPath string) error {
	resp, err := http.Get(payloadEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bomData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(absPath, bomFile), bomData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func currentPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}
	return filepath.Join(absPath, projectPath), nil
}
