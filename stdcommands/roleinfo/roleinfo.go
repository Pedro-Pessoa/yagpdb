package roleinfo

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/moderation"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var Command = &commands.TIDCommand{
	CmdCategory:  commands.CategoryGeneral,
	Name:         "Roleinfo",
	Aliases:      []string{"rinfo"},
	Description:  "Shows informations about the provided role",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Role", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		roleS := data.Args[0].Str()
		role := moderation.FindRole(data.GS, roleS)

		if role == nil {
			return "No role with the Name or ID `" + roleS + "` found", nil
		}

		embed := discordgo.MessageEmbed{
			Color: int(rand.Int63n(16777215)),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "ID",
					Value:  discordgo.StrID(role.ID),
					Inline: true,
				},
				{
					Name:   "Name",
					Value:  role.Name,
					Inline: true,
				},
				{
					Name:   "Color",
					Value:  strconv.Itoa(role.Color),
					Inline: true,
				},
				{
					Name:   "Mention",
					Value:  "`" + role.Mention() + "`",
					Inline: true,
				},
				{
					Name:   "Hoisted",
					Value:  strconv.FormatBool(role.Hoist),
					Inline: true,
				},
				{
					Name:   "Position",
					Value:  strconv.Itoa(role.Position),
					Inline: true,
				},
				{
					Name:   "Mentionable",
					Value:  strconv.FormatBool(role.Mentionable),
					Inline: true,
				},
				{
					Name:   "Managed",
					Value:  strconv.FormatBool(role.Managed),
					Inline: true,
				},
				{
					Name:   "Permissions",
					Value:  strings.Join(common.HumanizePermissions(role.Permissions), ", "),
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Created ",
			},
			Timestamp: (discordgo.SnowflakeTimestamp(role.ID)).Format(time.RFC3339),
		}

		return embed, nil
	},
}
