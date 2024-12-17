package testing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli-2/template/module"
)

func TestTemplateWrite(t *testing.T) {
	// Create a temporary directory to store the template and runtime files
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	templateFile, err := filepath.Abs("example_template.tmpl")
	if err != nil {
		t.Fatalf("Failed to get absolute path of template file: %v", err)
	}
	outputFile := filepath.Join(tempDir, "output.yml")

	customFields := map[string]string{
		"network":   "mainnet",
		"addresses": "addr1,addr2,addr3",
	}

	dataSource := &module.TemplateDataSource{
		ModuleConfigDir:     "/path/to/config",
		ModuleSecretFile:    "/path/to/secret.key",
		ModuleLogDir:        "/path/to/logs",
		ModuleJwtKeyFile:    "/path/to/jwt.key",
		HyperdriveDaemonUrl: "http://localhost:1234",
		CustomFields:        customFields,
	}

	tmpl := module.Template{
		Src: templateFile,
		Dst: outputFile,
	}

	err = tmpl.Write(dataSource)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedContent, err := os.ReadFile("expected_template_output.yml")
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	if string(outputContent) != string(expectedContent) {
		t.Errorf("Rendered output does not match expected.\nGot:\n%s\nExpected:\n%s", string(outputContent), string(expectedContent))
	}

}
