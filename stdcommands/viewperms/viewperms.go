package viewperms

import (
	"fmt"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryDebug,
	Name:        "ViewPerms",
	Description: "Shows you or the targets permissions in this channel",
	Arguments: []*dcmd.ArgDef{
		{Name: "target", Type: dcmd.UserID, Default: int64(0)},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		var target *dstate.MemberState

		if targetID := data.Args[0].Int64(); targetID == 0 {
			target = data.MS
		} else {
			var err error
			target, err = bot.GetMember(data.GS.ID, targetID)
			if err != nil {
				if common.IsDiscordErr(err, discordgo.ErrCodeUnknownMember) {
					return "Unknown member", nil
				}

				return nil, err
			}
		}

		perms, err := data.GS.MemberPermissionsMS(true, data.CS.ID, target)
		if err != nil {
			return "Unable to calculate perms", err
		}

		humanized := common.HumanizePermissions(int64(perms))
		return fmt.Sprintf("Perms of %s in this channel\n`%d`\n%s", target.Username, perms, strings.Join(humanized, ", ")), nil
	},
}
