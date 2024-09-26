package common

type DingDingNotify struct {
	Url           string `yaml:"url"`
	Secret        string `yaml:"secret"`
	Enabled       string `yaml:"enabled"`
	AlarmInterval int    `yaml:"alarm_interval"`
}

type DingDingReq struct {
	MsgType string             `json:"msgtype"`
	Text    DingDingReqContent `json:"text"`
	At      DingDingAt         `json:"at"`
}

type DingDingReqContent struct {
	Content string `json:"content"`
}

type DingDingAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	AsAtAll   bool     `json:"isAtAll"`
}