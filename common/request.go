package common

type RequiredReq struct {
	SignType    string `json:"signType,required" vd:"@:len($)>0; msg:'signType不能为空'"`
	Sign        string `json:"sign,required" vd:"@:len($)>0; msg:'sign不能为空'"`
	Filename    string `json:"filename,required" vd:"@:len($)>0; msg:'filename不能为空'"`
	EncryptType string `json:"encryptType,required" vd:"@:len($)>0; msg:'encryptType不能为空'"`
	BizContent  string `json:"bizContent,required" vd:"@:len($)>0; msg:'bizContent不能为空'"`
	Timestamp   string `json:"timestamp,required" vd:"@:len($)>0; msg:'timestamp不能为空'"`
	Version     string `json:"version,required" vd:"@:len($)>0; msg:'version不能为空'"`
}
