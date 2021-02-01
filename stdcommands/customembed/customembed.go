package customembed

import (
	"encoding/json"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var Command = &commands.TIDCommand{
	CmdCategory:         commands.CategoryFun,
	Name:                "CustomEmbed",
	Aliases:             []string{"ce"},
	Description:         "Creates an embed from what you give it in json form: https://docs.yagpdb.xyz/others/custom-embeds",
	LongDescription:     "Example: `-ce {\"title\": \"hello\", \"description\": \"wew\"}`",
	RequiredArgs:        1,
	RequireDiscordPerms: []int64{discordgo.PermissionManageMessages},
	Arguments: []*dcmd.ArgDef{
		{Name: "Json", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		var parsed *discordgo.MessageEmbed
		err := json.Unmarshal([]byte(data.Args[0].Str()), &parsed)
		if err != nil {
			return "Failed parsing json: " + err.Error(), err
		}

		return parsed, nil
	},
}
