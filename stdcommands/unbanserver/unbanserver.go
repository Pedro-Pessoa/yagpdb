package unbanserver

import (
	"github.com/mediocregopher/radix/v3"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "unbanserver",
	Description:          ";))",
	HideFromHelp:         true,
	RequiredArgs:         1,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "server", Type: dcmd.String},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		var unbanned bool
		err := common.RedisPool.Do(radix.Cmd(&unbanned, "SREM", "banned_servers", data.Args[0].Str()))
		if err != nil {
			return nil, err
		}

		if !unbanned {
			return "Server wasnt banned", nil
		}

		return "Unbanned server", nil
	}),
}
