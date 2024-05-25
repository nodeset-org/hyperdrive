package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/nodeset-org/hyperdrive-daemon/client"
	"github.com/nodeset-org/hyperdrive-daemon/internal/tests"
	"github.com/nodeset-org/hyperdrive-daemon/server"
	"github.com/nodeset-org/hyperdrive-daemon/shared/config"
)

// Various singleton variables used for testing
var (
	testMgr   *tests.TestManager    = nil
	wg        *sync.WaitGroup       = nil
	serverMgr *server.ServerManager = nil
	logger    *slog.Logger          = nil
	apiClient *client.ApiClient     = nil
)

// Initialize a common server used by all tests
func TestMain(m *testing.M) {
	wg = &sync.WaitGroup{}
	var err error
	testMgr, err = tests.NewTestManager()
	if err != nil {
		fail("error creating test manager: %v", err)
	}
	logger = testMgr.Logger

	// Create the server
	ip := "localhost"
	serverMgr, err = server.NewServerManager(testMgr.ServiceProvider, ip, 0, wg)
	if err != nil {
		fail("error creating server: %v", err)
	}

	// Create the client
	urlString := fmt.Sprintf("http://%s:%d/%s", ip, serverMgr.GetPort(), config.HyperdriveApiClientRoute)
	url, err := url.Parse(urlString)
	if err != nil {
		fail("error parsing client URL [%s]: %v", urlString, err)
	}
	apiClient = client.NewApiClient(url, logger, nil)

	// Run tests
	code := m.Run()

	// Revert to the baseline after testing is done
	cleanup()

	// Done
	os.Exit(code)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	cleanup()
	os.Exit(1)
}

func cleanup() {
	if serverMgr != nil {
		serverMgr.Stop()
		wg.Wait()
		logger.Info("Stopped server")
	}
	if testMgr != nil {
		testMgr.Cleanup()
		logger.Info("Cleaned up test manager")
	}
}
