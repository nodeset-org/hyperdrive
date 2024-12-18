package testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli-2/template/adapter"
)

func TestAdapterTemplateWrite(t *testing.T) {
	// Create a temporary directory to store the template and runtime files
	tempDir, err := os.MkdirTemp("", "adapter_template_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	templateFile, err := filepath.Abs("example_adapter_template.tmpl")
	if err != nil {
		t.Fatalf("Failed to get absolute path of template file: %v", err)
	}
	outputFile := filepath.Join(tempDir, "output.yml")

	dataSource := &adapter.AdapterDataSource{
		GetProjectName:   "hd2-service",
		ModuleConfigDir:  "/path/to/config",
		ModuleSecretFile: "/path/to/secret",
		ModuleLogDir:     "/path/to/log",
		ModuleJwtKeyFile: "/path/to/jwt",
	}

	tmpl := adapter.AdapterTemplate{
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

	expectedContent, err := os.ReadFile("expected_adapter_output.yml")
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	if strings.TrimSpace(string(outputContent)) != strings.TrimSpace(string(expectedContent)) {
		t.Errorf("Rendered output does not match expected.\nGot:\n%s\nExpected:\n%s", string(outputContent), string(expectedContent))
	}

}
