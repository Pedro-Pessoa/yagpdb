package dcallvoice

import (
	"fmt"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "dcallvoice",
	Description:          "Disconnects from all the voice channels the bot is in",
	HideFromHelp:         true,
	IsModCmd:             true,
	RunFunc: util.RequireBotAdmin(func(data *dcmd.Data) (interface{}, error) {
		vcs := make([]*discordgo.VoiceState, 0)
		guilds := bot.State.GuildsSlice(true)

		for _, g := range guilds {
			vc := g.VoiceState(true, common.BotUser.ID)
			if vc != nil {
				vcs = append(vcs, vc)
				go bot.ShardManager.SessionForGuild(g.ID).GatewayManager.ChannelVoiceLeave(g.ID)
			}
		}

		return fmt.Sprintf("Leaving %d voice channels...", len(vcs)), nil
	}),
}
