package undelete

import (
	"fmt"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var Command = &commands.TIDCommand{
	CmdCategory:  commands.CategoryTool,
	Name:         "Undelete",
	Aliases:      []string{"ud"},
	Description:  "Views your recent deleted messages, or all users deleted messages (with \"-a\" and manage messages perm) in this channel",
	RequiredArgs: 0,
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "a", Name: "all"},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		allUsers := data.Switch("a").Value != nil && data.Switch("a").Value.(bool)

		if allUsers {
			if ok, err := bot.AdminOrPermMS(data.CS.ID, data.MS, discordgo.PermissionManageMessages); !ok || err != nil {
				if err != nil {
					return nil, err
				} else if !ok {
					return "You need `Manage Messages` permissions to view all users deleted messages", nil
				}
			}
		}

		var resp strings.Builder
		resp.WriteString("Up to 10 last deleted messages (last hour or 12 hours for premium): \n\n")
		var numFound int

		data.GS.RLock()
		defer data.GS.RUnlock()

		for i := len(data.CS.Messages) - 1; i >= 0 && numFound < 10; i-- {
			msg := data.CS.Messages[i]

			if !msg.Deleted {
				continue
			}

			if !allUsers && msg.Author.ID != data.Msg.Author.ID {
				continue
			}

			var precision common.DurationFormatPrecision
			since := time.Since(msg.ParsedCreated)

			switch {
			case since < time.Minute:
				precision = common.DurationPrecisionSeconds
			case since < time.Hour:
				precision = common.DurationPrecisionMinutes
			default:
				precision = common.DurationPrecisionHours
			}

			// Match found!
			timeSince := common.HumanizeDuration(precision, since)

			resp.WriteString(fmt.Sprintf("`%s ago (%s)` **%s**#%s: %s\n\n", timeSince, msg.ParsedCreated.UTC().Format(time.ANSIC), msg.Author.Username, msg.Author.Discriminator, msg.ContentWithMentionsReplaced()))
			numFound++
		}

		if numFound == 0 {
			resp.WriteString("none...")
		}

		return resp.String(), nil
	},
}
