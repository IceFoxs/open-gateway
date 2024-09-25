package dingding

import (
	"github.com/IceFoxs/open-gateway/common"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	"net/url"
	"strconv"
	"testing"
)

func TestDingdingSign(t *testing.T) {
	timestamp := 1727277749036
	secret := "12313131"
	signStr := strconv.FormatInt(int64(timestamp), 10) + "\n" + secret
	sign := rsaUtil.HmacSha256ToBase64(secret, signStr)
	t.Log(sign)
	t.Log(url.QueryEscape(sign))
	ddc, _ := NewDingDingNotifyClient(&common.DingDingNotify{
		Url:    "https://oapi.dingtalk.com/robot/send?access_token=ad6f72e51f4524a14e670ec8a3cabe82d38ec6c46295d29935929d96c0cfbf29",
		Secret: "SECba5111b9730d21ea00e6ca6cfa05bc34d55296a93a4087b113a4cf4673f9bcec",
	})
	t.Log(ddc.generateUrl())
	_, err := ddc.Req("自研支付成员您好！！！", false)
	if err != nil {
		t.Log(err.Error())
	}
}
