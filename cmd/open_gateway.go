package cmd

import (
	_ "github.com/IceFoxs/open-gateway/cmd/imports"
	"github.com/IceFoxs/open-gateway/server"
	_ "github.com/IceFoxs/open-gateway/sync/imports"
)

import (
	"github.com/IceFoxs/open-gateway/server/router"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Start() {
	h, err := server.CreateServer()
	if err != nil {
		hlog.Errorf("create server failed: %s", err.Error())
	}
	router.AddRouter(h)
	h.Spin()
}
