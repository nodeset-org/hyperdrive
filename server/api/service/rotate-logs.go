package service

import (
	"errors"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/rocket-pool/node-manager-core/api/server"
	"github.com/rocket-pool/node-manager-core/api/types"
)

// ===============
// === Factory ===
// ===============

type serviceRotateLogsContextFactory struct {
	handler *ServiceHandler
}

func (f *serviceRotateLogsContextFactory) Create(args url.Values) (*serviceRotateLogsContext, error) {
	c := &serviceRotateLogsContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *serviceRotateLogsContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*serviceRotateLogsContext, types.SuccessData](
		router, "rotate", f, f.handler.logger.Logger, f.handler.serviceProvider.ServiceProvider,
	)
}

// ===============
// === Context ===
// ===============

type serviceRotateLogsContext struct {
	handler *ServiceHandler
}

func (c *serviceRotateLogsContext) PrepareData(data *types.SuccessData, opts *bind.TransactOpts) (types.ResponseStatus, error) {
	sp := c.handler.serviceProvider
	apiLog := sp.GetApiLogger()
	tasksLog := sp.GetTasksLogger()

	err := errors.Join(
		apiLog.Rotate(),
		tasksLog.Rotate(),
	)
	if err != nil {
		return types.ResponseStatus_Error, err
	}
	return types.ResponseStatus_Success, nil
}
