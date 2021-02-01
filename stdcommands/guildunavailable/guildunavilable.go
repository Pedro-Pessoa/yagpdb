package guildunavailable

import (
	"fmt"

	"github.com/Pedro-Pessoa/tidbot/bot/botrest"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory:  commands.CategoryDebug,
	Name:         "IsGuildUnavailable",
	Description:  "Returns wether the specified guild is unavilable or not",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "guildid", Type: dcmd.Int, Default: int64(0)},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		gID := data.Args[0].Int64()
		guild, err := botrest.GetGuild(gID)
		if err != nil {
			return "Uh oh", err
		}

		return fmt.Sprintf("Guild (%d) unavilable: %v", guild.ID, guild.Unavailable), nil
	},
}
