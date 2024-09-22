package sync

import (
	"github.com/IceFoxs/open-gateway/common"
	"github.com/IceFoxs/open-gateway/sync/config"
	"sync"
)

var (
	confChangeClientHelper *ConfChangeClientHelper
	once                   sync.Once
)

type ConfChangeClientHelper struct {
	confChangeClient *config.ConfChangeClient
}

func GetConfChangeClientHelper() *ConfChangeClientHelper {
	once.Do(initConfChangeClientHelper)
	return confChangeClientHelper
}

func initConfChangeClientHelper() {
	confChangeClientHelper = &ConfChangeClientHelper{}
}
func (c *ConfChangeClientHelper) Publish(confType, confGroup, confContent string) {
	cc := c.confChangeClient
	if cc != nil {
		(*cc).Publish(confType, confGroup, confContent)
	}
}

func (c *ConfChangeClientHelper) Subscribe(confType, confGroup string, listener common.Listener) {
	cc := c.confChangeClient
	if cc != nil {
		(*cc).Subscribe(confType, confGroup, listener)
	}
}
func (c *ConfChangeClientHelper) BuildConfChangeClient(confChangeClient *config.ConfChangeClient) {
	c.confChangeClient = confChangeClient
}
