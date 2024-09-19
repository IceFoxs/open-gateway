package nacos_client_test

import (
	"context"
	"fmt"
	c "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"testing"
)

func TestPostJson(t *testing.T) {
	client, err := c.NewClient()
	if err != nil {
		panic(err)
	}
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI("http://127.0.0.1:8888/api/json")
	req.SetOptions(config.WithSD(false))
	req.SetBody([]byte("{\"bizContent\":\"{\\\"pageNo\\\":\\\"1\\\",\\\"pageSize\\\":10}\",\"encryptType\":\"NONE\",\"filename\":\"FPS_MANAGER_GETBANKCHANNELCONFIG_REQ_990102_20240813164335632.json\",\"sign\":\"\",\"signType\":\"NONE\",\"timestamp\":\"2024-08-13T16:43:35\",\"version\":\"2.0\"}"))
	err = client.Do(context.Background(), req, res)
	if err != nil {
		return
	}
	fmt.Printf("code=%d,body=%s\n", res.StatusCode(), string(res.Body()))
}
