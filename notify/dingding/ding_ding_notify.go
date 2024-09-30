package dingding

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/IceFoxs/open-gateway/common"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/dubbogo/gost/log/logger"
	"net/url"
	"strconv"
	"time"
)

type DingDingNotifyClient struct {
	dc     *common.DingDingNotify
	Client *client.Client
}

func NewDingDingNotifyClient(dc *common.DingDingNotify) (*DingDingNotifyClient, error) {
	clientCfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := client.NewClient(
		client.WithTLSConfig(clientCfg),
		client.WithDialer(standard.NewDialer()),
	)
	if err != nil {
		return nil, err
	}
	return &DingDingNotifyClient{
		dc:     dc,
		Client: c,
	}, nil
}
func (d *DingDingNotifyClient) Req(dingdingReq common.DingDingReq) (bool, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI(d.generateUrl())
	req.SetOptions(config.WithSD(false))
	bt, _ := json.Marshal(dingdingReq)
	req.SetBody(bt)
	err := d.Client.Do(context.Background(), req, res)
	if err != nil {
		return false, err
	}
	logger.Infof("%s", string(res.Body()))
	if res.StatusCode() == 200 {
		return true, nil
	}
	return false, err
}
func (d *DingDingNotifyClient) ReqMarkdown(title, param string, atAll bool) (bool, error) {
	ddReq := common.DingDingReq{
		MsgType:  "markdown",
		Markdown: common.Markdown{Title: title, Text: param},
		At: common.DingDingAt{
			AsAtAll:   atAll,
			AtMobiles: []string{},
			AtUserIds: []string{},
		},
	}
	return d.Req(ddReq)
}
func (d *DingDingNotifyClient) ReqText(param string, atAll bool) (bool, error) {
	ddReq := common.DingDingReq{
		MsgType: "text",
		Text:    common.Markdown{Content: param},
		At: common.DingDingAt{
			AsAtAll:   atAll,
			AtMobiles: []string{},
			AtUserIds: []string{},
		},
	}
	return d.Req(ddReq)
}
func (d *DingDingNotifyClient) generateUrl() string {
	timestamp := time.Now().UnixNano() / 1e6
	signStr := strconv.FormatInt(timestamp, 10) + "\n" + d.dc.Secret
	sign := rsaUtil.HmacSha256ToBase64(d.dc.Secret, signStr)
	return d.dc.Url + "&timestamp=" + strconv.FormatInt(timestamp, 10) + "&sign=" + url.QueryEscape(sign)
}
