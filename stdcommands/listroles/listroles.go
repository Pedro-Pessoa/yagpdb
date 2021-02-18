package listroles

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dutil"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryTool,
	Name:        "ListRoles",
	Aliases:     []string{"lr", "ur"},
	Description: "List roles, their id's, color hex code, and 'mention everyone' perms (useful if you wanna double check to make sure you didn't give anyone mention everyone perms that shouldn't have it)",
	Arguments: []*dcmd.ArgDef{
		{Name: "User", Type: &commands.MemberArg{}},
	},

	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "nomanaged", Name: "Don't list managed/bot roles"},
	},

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		var out strings.Builder
		var noMana bool
		var member *dstate.MemberState

		if data.Args[0].Value != nil {
			member = data.Args[0].Value.(*dstate.MemberState)
		}

		if data.Switches["nomanaged"].Value != nil && data.Switches["nomanaged"].Value.(bool) && member == nil {
			noMana = true
		}

		data.GS.Lock()
		defer data.GS.Unlock()

		sort.Sort(dutil.Roles(data.GS.Guild.Roles))

		var counter int
		if member != nil {
			if len(member.Roles) > 0 {
				for _, roleID := range member.Roles {
					for _, r := range data.GS.Guild.Roles {
						if roleID == r.ID {
							counter++
							me := r.Permissions&discordgo.PermissionAdministrator != 0 || r.Permissions&discordgo.PermissionMentionEveryone != 0
							out.WriteString(fmt.Sprintf("`%-25s: %-19s #%-6x  ME:%5t`\n", r.Name, strconv.FormatInt(r.ID, 10), r.Color, me))
						}
					}
				}
			} else {
				return "Esse membro não tem cargos.", nil
			}
		} else {
			for _, r := range data.GS.Guild.Roles {
				if noMana && r.Managed {
					continue
				} else {
					counter++
					me := r.Permissions&discordgo.PermissionAdministrator != 0 || r.Permissions&discordgo.PermissionMentionEveryone != 0
					out.WriteString(fmt.Sprintf("`%-25s: %-19s #%-6x  ME:%5t`\n", r.Name, strconv.FormatInt(r.ID, 10), r.Color, me))
				}
			}
		}

		return ("Total role count: " + strconv.Itoa(counter) + "\n(ME = mention everyone perms)\n" + out.String()), nil
	},
}
