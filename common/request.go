package common

import (
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
)

const (
	REQ          = "req"
	FILENAME_REQ = "filenameReq"
	REQ_BODY     = "reqBody"
)

type RequiredReq struct {
	SignType    string `json:"signType,required" vd:"@:len($)>0; msg:'signType不能为空'"`
	Sign        string `json:"sign,required"`
	Filename    string `json:"filename,required" vd:"@:len($)>0; msg:'filename不能为空'"`
	EncryptType string `json:"encryptType,required" vd:"@:len($)>0; msg:'encryptType不能为空'"`
	BizContent  string `json:"bizContent,required"`
	Timestamp   string `json:"timestamp,required" vd:"@:len($)>0; msg:'timestamp不能为空'"`
	Version     string `json:"version,required" vd:"@:len($)>0; msg:'version不能为空'"`
}

type DecryptContentReq struct {
	AppId          string `json:"appId,required"`
	EncryptContent string `json:"encryptContent,required"`
}

type CommonRes struct {
	BizContent string `json:"bizContent"`
	ErrorMsg   string `json:"errorMsg"`
	Sign       string `json:"sign"`
	StatusCode int    `json:"statusCode"`
}

func Error(code int, msg string) CommonRes {
	hlog.Errorf("response : %s ,%s", strconv.Itoa(code), msg)
	return CommonRes{
		Sign:       "NONE",
		ErrorMsg:   msg,
		StatusCode: code,
	}
}

func ErrorWithSign(code int, msg string, sign string) CommonRes {
	hlog.Errorf("response code :%s ,msg: %s", strconv.Itoa(code), msg)
	return CommonRes{
		Sign:       sign,
		ErrorMsg:   msg,
		StatusCode: code,
	}
}

func Success(code int, bizContent string, sign string) CommonRes {
	hlog.Infof("response code :%s ,content:%s,sign:%s", strconv.Itoa(code), bizContent, sign)
	return CommonRes{
		ErrorMsg:   "请求成功",
		StatusCode: code,
		BizContent: bizContent,
		Sign:       sign,
	}
}

func Succ(code int, msg string, i interface{}, sign string) CommonRes {
	content, err := json.Marshal(i)
	if err != nil {
		return Error(500, err.Error())
	}
	hlog.Infof("response code: %s ,content: %s ,sign: %s", strconv.Itoa(code), content, sign)
	return CommonRes{
		ErrorMsg:   msg,
		StatusCode: code,
		BizContent: string(content),
		Sign:       sign,
	}
}

func SuccContent(code int, msg string, content string, sign string) CommonRes {
	hlog.Infof("response code: %s,msg:%s ,content: %s ,sign: %s", strconv.Itoa(code), msg, content, sign)
	return CommonRes{
		ErrorMsg:   msg,
		StatusCode: code,
		BizContent: content,
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
