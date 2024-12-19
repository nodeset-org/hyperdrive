package testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli-2/template/metadata"
)

func TestMetadataTemplateWrite(t *testing.T) {
	templateFile, err := filepath.Abs("example_metadata_template.tmpl")
	if err != nil {
		t.Fatalf("Failed to get absolute path of template file: %v", err)
	}
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedContent, err := os.ReadFile("expected_metadata_output.yml")
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	customFields := map[string]string{
		"network":   "mainnet",
		"addresses": "addr1,addr2,addr3",
	}

	dataSource := &metadata.MetadataDataSource{
		CustomFields: customFields,
	}

	tmpl := metadata.MetadataTemplate{
		Content: string(templateContent),
	}

	output, err := tmpl.Write(dataSource)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if strings.TrimSpace(string(output)) != strings.TrimSpace(string(expectedContent)) {
		t.Errorf("Rendered output does not match expected.\nGot:\n%s\nExpected:\n%s", string(output), string(expectedContent))
	}

}
