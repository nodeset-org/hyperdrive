package tests

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"testing"

	"github.com/nodeset-org/osha"
	"github.com/rocket-pool/node-manager-core/log"
)

// Various singleton variables used for testing
var (
	testMgr *osha.TestManager = nil
	logger  *slog.Logger      = nil
)

// Initialize a common server used by all tests
func TestMain(m *testing.M) {
	var err error

	// Create a new test manager
	testMgr, err = osha.NewTestManager()
	if err != nil {
		fail("error creating test manager: %v", err)
	}
	logger = testMgr.GetLogger()

	// Run tests
	code := m.Run()

	// Clean up and exit
	mainCleanup()
	os.Exit(code)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	mainCleanup()
	os.Exit(1)
}

func mainCleanup() {
	if testMgr == nil {
		return
	}
	err := testMgr.Close()
	if err != nil {
		logger.Error("Error closing test manager", log.Err(err))
	}
	testMgr = nil
}

func handle_panics() {
	// Handle panics
	r := recover()
	if r != nil {
		debug.PrintStack()
		fail("Recovered from panic: %v", r)
	}
}
