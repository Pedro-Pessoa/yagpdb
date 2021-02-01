package youtube

import (
	"fmt"

	"github.com/mediocregopher/radix/v3"

	"github.com/Pedro-Pessoa/tidbot/common"
)

func (p *Plugin) Status() (string, string) {
	var unique int
	_ = common.RedisPool.Do(radix.Cmd(&unique, "ZCARD", "youtube_subbed_channels"))

	var numChannels int
	common.GORM.Model(&ChannelSubscription{}).Count(&numChannels)

	return "Youtube", fmt.Sprintf("%d/%d", unique, numChannels)
}
