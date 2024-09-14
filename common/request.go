package common

import "github.com/cloudwego/hertz/pkg/common/json"

type RequiredReq struct {
	SignType    string `json:"signType,required" vd:"@:len($)>0; msg:'signType不能为空'"`
	Sign        string `json:"sign,required" vd:"@:len($)>0; msg:'sign不能为空'"`
	Filename    string `json:"filename,required" vd:"@:len($)>0; msg:'filename不能为空'"`
	EncryptType string `json:"encryptType,required" vd:"@:len($)>0; msg:'encryptType不能为空'"`
	BizContent  string `json:"bizContent,required" vd:"@:len($)>0; msg:'bizContent不能为空'"`
	Timestamp   string `json:"timestamp,required" vd:"@:len($)>0; msg:'timestamp不能为空'"`
	Version     string `json:"version,required" vd:"@:len($)>0; msg:'version不能为空'"`
}

type CommonRes struct {
	BizContent string `json:"bizContent"`
	ErrorMsg   string `json:"errorMsg"`
	Sign       string `json:"sign"`
	StatusCode int    `json:"statusCode"`
}

func Error(code int, msg string) CommonRes {
	return CommonRes{
		Sign:       "NONE",
		ErrorMsg:   msg,
		StatusCode: code,
	}
}

func Success(code int, bizContent string, sign string) CommonRes {
	return CommonRes{
		ErrorMsg:   "请求成功",
		StatusCode: code,
		BizContent: bizContent,
		Sign:       sign,
	}
}

func Succ(code int, i interface{}, sign string) CommonRes {
	content, err := json.Marshal(i)
	if err != nil {
		return Error(500, err.Error())
	}
	return CommonRes{
		ErrorMsg:   "请求成功",
		StatusCode: code,
		BizContent: string(content),
		Sign:       sign,
	}
}

func ToJSON(i interface{}) string {
	content, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(content)
}
