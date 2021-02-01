package setstatus

import (
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "setstatus",
	Description:          "Sets the bot's status and streaming url",
	HideFromHelp:         true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "status", Type: dcmd.String, Default: ""},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "url", Name: "streaming url", Type: dcmd.String, Default: ""},
	},
	RunFunc: util.RequireBotAdmin(func(data *dcmd.Data) (interface{}, error) {
		streamingURL := data.Switch("url").Str()
		bot.SetStatus(streamingURL, data.Args[0].Str())
		return "Doneso", nil
	}),
}
