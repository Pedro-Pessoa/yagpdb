package editrole

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/moderation"
	"github.com/jonas747/yagpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory:     commands.CategoryTool,
	Name:            "EditRole",
	Aliases:         []string{"ERole"},
	Description:     "Edits a role",
	LongDescription: "Requires the manage roles permission and the bot and your highest role being above the edited role. Role permissions follow discord standard encoding can can be calculated [here](https://discordapp.com/developers/docs/topics/permissions)",
	RequiredArgs:    1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Role", Type: dcmd.String},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "name", Help: "Role name - String", Type: dcmd.String, Default: ""},
		{Switch: "color", Help: "Role color - Either hex code or name", Type: dcmd.String, Default: ""},
		{Switch: "mention", Help: "Role Mentionable - 1 for true 0 for false", Type: &dcmd.IntArg{Min: 0, Max: 1}},
		{Switch: "hoist", Help: "Role Hoisted - 1 for true 0 for false", Type: &dcmd.IntArg{Min: 0, Max: 1}},
		{Switch: "perms", Help: "Role Permissions - 0 to 2147483647", Type: &dcmd.IntArg{Min: 0, Max: 2147483647}},
	},
	RunFunc:             cmdFuncEditRole,
	GuildScopeCooldown:  15,
	RequireDiscordPerms: []int64{discordgo.PermissionManageRoles},
}

func cmdFuncEditRole(data *dcmd.Data) (interface{}, error) {
	if util.IsExecedByCC(data) && !bot.IsGuildWhiteListed(data.GS.ID) {
		return "", errors.New("Esse comando não pode ser invocado através de um Custom Command.")
	}

	roleS := data.Args[0].Str()
	role := moderation.FindRole(data.GS, roleS)

	if role == nil {
		return "No role with the Name or ID`" + roleS + "` found", nil
	}

	data.GS.RLock()
	if !bot.IsMemberAboveRole(data.GS, data.MS, role) {
		data.GS.RUnlock()
		return "Can't edit roles above you", nil
	}
	data.GS.RUnlock()

	change := false
	cID := data.CS.ID

	name := role.Name
	if n := data.Switch("name").Str(); n != "" {
		name = limitString(n, 100)
		change = true
	}

	color := role.Color
	if c := data.Switch("color").Str(); c != "" {
		parsedColor, ok := util.ParseColor(c)
		if !ok {
			return "Unknown color: " + c + ", can be either hex color code or name for a known color", nil
		}
		color = parsedColor
		change = true
	}

	mentionable := role.Mentionable
	if m := data.Switch("mention"); m != nil {
		mentionable = m.Bool()
		change = true
	}

	hoisted := role.Hoist
	if h := data.Switch("hoist"); h != nil {
		hoisted = h.Bool()
		change = true
	}

	perms := role.Permissions
	if p := data.Switch("perms"); p != nil {
		perms = p.Int()
		change = true
	}

	if change {
		_, err := common.BotSession.GuildRoleEdit(data.GS.ID, role.ID, name, color, hoisted, perms, mentionable)
		if err != nil {
			return nil, err
		}
	}

	_, err := common.BotSession.ChannelMessageSendComplex(cID, &discordgo.MessageSend{
		Content:         fmt.Sprintf("__**Edited Role (%d) properties to :**__\n\n**Name **: `%s`\n**Color **: `%d`\n**Mentionable **: `%t`\n**Hoisted **: `%t`\n**Permissions **: `%d`", role.ID, name, color, mentionable, hoisted, perms),
		AllowedMentions: discordgo.AllowedMentions{},
	})

	if err != nil {
		return nil, err
	}

	return nil, err
}

// limitstring cuts off a string at max l length, supports multi byte characters
func limitString(s string, l int) string {
	if len(s) <= l {
		return s
	}

	lastValidLoc := 0
	for i := range s {
		if i > l {
			break
		}
		lastValidLoc = i
	}

	return s[:lastValidLoc]
}
