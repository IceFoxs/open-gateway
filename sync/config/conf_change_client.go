package config

import "github.com/IceFoxs/open-gateway/common"

type ConfChangeClient interface {
	Publish(confType string, confGroup string, confContent string)
	Subscribe(confType string, confGroup string, listener common.Listener)
}
