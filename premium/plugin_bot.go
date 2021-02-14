package premium

import (
	"time"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

var _ bot.BotInitHandler = (*Plugin)(nil)
var _ commands.CommandProvider = (*Plugin)(nil)

func (p *Plugin) BotInit() {
	bot.State.CustomLimitProvider = p
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, cmdGenerateCode, cmdDeleteCode, cmdPremiumInfo)
}

const (
	NormalStateMaxMessages   = 1000
	NormalStateMaxMessageAge = time.Hour

	PremiumStateMaxMessags    = 10000
	PremiumStateMaxMessageAge = time.Hour * 12
)

func (p *Plugin) MessageLimits(gs *dstate.GuildState) (maxMessages int, maxMessageAge time.Duration) {
	if gs == nil {
		return NormalStateMaxMessages, NormalStateMaxMessageAge
	}

	premium, err := IsGuildPremiumCached(gs.ID)
	if err != nil {
		logger.WithError(err).WithField("guild", gs.ID).Error("Failed checking if guild is premium")
	}

	if premium {
		return PremiumStateMaxMessags, PremiumStateMaxMessageAge
	}

	return NormalStateMaxMessages, NormalStateMaxMessageAge
}
