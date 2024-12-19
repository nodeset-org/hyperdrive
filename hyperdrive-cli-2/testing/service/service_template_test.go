package testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nodeset-org/hyperdrive/hyperdrive-cli-2/template/service"
)

func TestServiceTemplateWrite(t *testing.T) {
	// Create a temporary directory to store the template and runtime files
	tempDir, err := os.MkdirTemp("", "service_template_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	templateFile, err := filepath.Abs("example_service_template.tmpl")
	if err != nil {
		t.Fatalf("Failed to get absolute path of template file: %v", err)
	}
	outputFile := filepath.Join(tempDir, "output.yml")

	customFields := map[string]string{
		"network":   "mainnet",
		"addresses": "addr1,addr2,addr3",
	}

	dataSource := &service.ServiceDataSource{
		GetProjectName:       "hd2-service",
		ModuleConfigDir:      "/path/to/config",
		ModuleDataDir:        "/path/to/data",
		HyperdriveDaemonUrl:  "http://localhost:1234",
		HyperdriveJwtKeyFile: "/path/to/jwt",

		CustomFields: customFields,
	}

	tmpl := service.ServiceTemplate{
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

	expectedContent, err := os.ReadFile("expected_service_output.yml")
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}

	if strings.TrimSpace(string(outputContent)) != strings.TrimSpace(string(expectedContent)) {
		t.Errorf("Rendered output does not match expected.\nGot:\n%s\nExpected:\n%s", string(outputContent), string(expectedContent))
	}

}