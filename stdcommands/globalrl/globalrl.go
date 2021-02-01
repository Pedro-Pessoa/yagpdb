package globalrl

import (
	"time"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	Name:                 "globalrl",
	Description:          "Tests the global ratelimit functionality",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		rlEvt := &discordgo.RateLimit{
			URL: "Wew",
			TooManyRequests: &discordgo.TooManyRequests{
				Bucket:     "wewsss",
				Message:    "Too many!",
				RetryAfter: time.Duration(5000),
			},
		}

		go common.BotSession.HandleEvent("__RATE_LIMIT__", rlEvt)

		return "Done", nil
	}),
}
