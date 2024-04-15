package types

type NodesetStatus string

const (
	NodesetStatus_Unknown               NodesetStatus = ""
	NodesetStatus_RegisteredToStakewise NodesetStatus = "RegisteredToStakewise"
	NodesetStatus_UploadedStakewise     NodesetStatus = "UploadedStakewise"
	NodesetStatus_UploadedToNodeset     NodesetStatus = "UploadedToNodeset"
	NodesetStatus_Generated             NodesetStatus = "Generated"
)
