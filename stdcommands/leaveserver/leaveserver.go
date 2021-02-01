package leaveserver

import (
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "leaveserver",
	Description:          ";))",
	HideFromHelp:         true,
	RequiredArgs:         1,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "server", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		err := common.BotSession.GuildLeave(data.Args[0].Int64())
		if err == nil {
			return "Left " + data.Args[0].Str(), nil
		}

		return err, err
	}),
}
