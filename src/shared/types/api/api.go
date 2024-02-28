package api

import (
	"github.com/nodeset-org/eth-utils/eth"
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
