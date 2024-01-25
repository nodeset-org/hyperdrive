package example

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/shared/config"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

const (
	apiContainerName string = "api"
	binaryPath       string = "/go/bin/hyperdrive"
)

// ===============
// === Factory ===
// ===============

type uploadDepositDataContextFactory struct {
	handler *NodeHandler
}

func (f *uploadDepositDataContextFactory) Create(args url.Values) (*uploadDepositDataContext, error) {
	c := &uploadDepositDataContext{
		handler: f.handler,
	}

	// Check for required input args
	var msg string
	inputErrs := []error{
		server.GetStringFromValues("message", args, &msg),
	}

	// Decode the message
	var err error
	c.message, err = base64.StdEncoding.DecodeString(msg)
	inputErrs = append(inputErrs, err)
	return c, errors.Join(inputErrs...)
}

func (f *uploadDepositDataContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*uploadDepositDataContext, api.UploadDepositDataData](
		router, "upload-deposit-data", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type uploadDepositDataContext struct {
	handler *NodeHandler
	cfg     *config.HyperdriveConfig

	message []byte
}

func (c *uploadDepositDataContext) PrepareData(data *api.UploadDepositDataData) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Sign the message
	signedMessage, err := w.SignMessage(c.message)
	if err != nil {
		return fmt.Errorf("error signing message: %w", err)
	}

	// TODO: POST
	fmt.Sprint(signedMessage)

	return nil
}
