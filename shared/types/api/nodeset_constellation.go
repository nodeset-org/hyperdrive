package api

type NodeSetConstellation_GetRegistrationSignatureData struct {
	NotRegistered bool   `json:"notRegistered"`
	NotAuthorized bool   `json:"notAuthorized"`
	Signature     []byte `json:"signature"`
}

type NodeSetConstellation_GetAvailableMinipoolCount struct {
	NotRegistered bool `json:"notRegistered"`
	Count         int  `json:"count"`
}

type NodeSetConstellation_GetDepositSignatureData struct {
	NotRegistered      bool   `json:"notRegistered"`
	NotAuthorized      bool   `json:"notAuthorized"`
	LimitReached       bool   `json:"limitReached"`
	MissingExitMessage bool   `json:"missingExitMessage"`
	Signature          []byte `json:"signature"`
}
