package regex_test

import (
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/common/regex"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

func TestFilename(*testing.T) {
	defer func() {
		if err := recover(); err != nil {
			hlog.Errorf("filename error :%s", err.(error).Error())
		}
	}()
	name, err := regex.MatchFileName("FPS_MANAGER_CONFREFRESH_REQ_niyang_20240813164335632.json")
	if err != nil {
		hlog.Errorf("filename error :%s", err.Error())
		return
	}
	hlog.Infof("filename:%s", common.ToJSON(name))
	name, err = regex.MatchFileName("FPS_MANAGER_CONFREFRESH_1_REQ_niyang_20240813164335632.json")
	if err != nil {
		hlog.Errorf("filename error :%s", err.Error())
		return
	}
	hlog.Infof("filename:%s", common.ToJSON(name))

}
