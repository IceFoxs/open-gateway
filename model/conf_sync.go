package model

type ConfTypeRequest struct {
	ConfType    string `json:"confType" vd:"@:len($)>0; msg:'confType不能为空'"`
	ConfContent string `json:"confContent"`
}
