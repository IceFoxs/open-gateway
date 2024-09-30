package common

type DingDingNotify struct {
	Url           string `yaml:"url"`
	Secret        string `yaml:"secret"`
	Enabled       string `yaml:"enabled"`
	AlarmInterval int    `yaml:"alarm_interval"`
}

type DingDingReq struct {
	MsgType  string     `json:"msgtype"`
	Markdown Markdown   `json:"markdown"`
	Text     Markdown   `json:"text"`
	At       DingDingAt `json:"at"`
}

type Markdown struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Text    string `json:"text"`
}

type DingDingAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	AsAtAll   bool     `json:"isAtAll"`
}
