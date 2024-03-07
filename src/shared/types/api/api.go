package api

import (
	"github.com/rocket-pool/node-manager-core/eth"
)

type ApiResponse[Data any] struct {
	Data *Data `json:"data"`
}

type SuccessData struct {
}

type DataBatch[DataType any] struct {
	Batch []DataType `json:"batch"`
}

type TxInfoData struct {
	TxInfo *eth.TransactionInfo `json:"txInfo"`
}

type BatchTxInfoData struct {
	TxInfos []*eth.TransactionInfo `json:"txInfos"`
}
