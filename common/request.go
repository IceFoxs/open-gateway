package common

type RequiredReq struct {
	SignType    string `json:"signType,required"`
	Sign        string `json:"sign,required"`
	Filename    string `json:"filename,required" `
	EncryptType string `json:"encryptType,required"`
	BizContent  string `json:"bizContent,required"`
	Timestamp   string `json:"timestamp,required" `
	Version     string `json:"version,required"`
}
