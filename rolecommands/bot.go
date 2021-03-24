package rolecommands

import (
	"context"
	"database/sql"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Pedro-Pessoa/tidbot/analytics"
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/pubsub"
	"github.com/Pedro-Pessoa/tidbot/common/scheduledevents2"
	schEvtsModels "github.com/Pedro-Pessoa/tidbot/common/scheduledevents2/models"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/rolecommands/models"
)

func (p *Plugin) AddCommands() {
	const msgIDDocs = "Para conseguir o ID de uma mensagem você deve ativar o modo desenvolvidor no discord. Feito isso, clique com o botão direito na mensagem e depois em copiar ID."

	categoryRoleMenu := &dcmd.Category{
		Name:        "Rolemenu",
		Description: "Rolemenu commands",
		HelpEmoji:   "🔘",
		EmbedColor:  0x42b9f4,
	}

	commands.AddRootCommands(p,
		&commands.TIDCommand{
			CmdCategory: commands.CategoryTool,
			Name:        "Role",
			Aliases:     []string{"cargo"},
			Description: "Alterna um cargo em você ou lista todos os cargos disponíveis. Os cargos tem que ser configurados no painel de controle. ",
			Arguments: []*dcmd.ArgDef{
				{Name: "Role", Type: dcmd.String},
			},
			RunFunc: CmdFuncRole,
		})

	cmdCreate := &commands.TIDCommand{
		Name:                "Create",
		CmdCategory:         categoryRoleMenu,
		Aliases:             []string{"c"},
		Description:         "Set up a role menu.",
		LongDescription:     "Specify a message with -m to use an existing message instead of having the bot make one\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Group", Type: dcmd.String},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "m", Name: "Message ID", Type: &dcmd.IntArg{}},
			{Switch: "nodm", Name: "Disable DM"},
			{Switch: "rr", Name: "Remove role on reaction removed"},
			{Switch: "skip", Name: "Number of roles to skip", Default: 0, Type: dcmd.Int},
		},
		RunFunc: cmdFuncRoleMenuCreate,
	}

	cmdRemoveRoleMenu := &commands.TIDCommand{
		Name:                "Remove",
		CmdCategory:         categoryRoleMenu,
		Description:         "Removes a rolemenu from a message.",
		LongDescription:     "The message won't be deleted and the bot will not do anything with reactions on that message\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Message ID", Type: dcmd.Int},
		},
		RunFunc: cmdFuncRoleMenuRemove,
	}

	cmdUpdate := &commands.TIDCommand{
		Name:                "Update",
		CmdCategory:         categoryRoleMenu,
		Aliases:             []string{"u"},
		Description:         "Updates a rolemenu, toggling the provided flags and adding missing options, aswell as updating the order.",
		LongDescription:     "\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Message ID", Type: dcmd.Int},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "nodm", Name: "Disable DM"},
			{Switch: "rr", Name: "Remove role on reaction removed"},
		},
		RunFunc: cmdFuncRoleMenuUpdate,
	}

	cmdResetReactions := &commands.TIDCommand{
		Name:                "ResetReactions",
		CmdCategory:         categoryRoleMenu,
		Aliases:             []string{"reset"},
		Description:         "Removes all reactions on the specified menu message and re-adds them.",
		LongDescription:     "Can be used to fix the order after updating it.\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Message ID", Type: dcmd.Int},
		},
		RunFunc: cmdFuncRoleMenuResetReactions,
	}

	cmdEditOption := &commands.TIDCommand{
		Name:                "EditOption",
		CmdCategory:         categoryRoleMenu,
		Aliases:             []string{"edit"},
		Description:         "Allows you to reassign the emoji of an option, tip: use ResetReactions afterwards.",
		LongDescription:     "\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Message ID", Type: dcmd.Int},
		},
		RunFunc: cmdFuncRoleMenuEditOption,
	}

	cmdFinishSetup := &commands.TIDCommand{
		Name:                "Complete",
		CmdCategory:         categoryRoleMenu,
		Aliases:             []string{"finish"},
		Description:         "Marks the menu as done.",
		LongDescription:     "\n\n" + msgIDDocs,
		RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
		RequiredArgs:        1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Message ID", Type: dcmd.Int},
		},
		RunFunc: cmdFuncRoleMenuComplete,
	}

	menuContainer := commands.CommandSystem.Root.Sub("RoleMenu", "rmenu", "rm")

	const notFoundMessage = "Unknown rolemenu command, if you've used this before it was recently revamped.\nTry almost the same command but `rolemenu create ...` and `rolemenu update ...` instead (replace '...' with the rest of the command).\nSee `help rolemenu` for all rolemenu commands."
	menuContainer.NotFound = commands.CommonContainerNotFoundHandler(menuContainer, notFoundMessage)

	menuContainer.AddCommand(cmdCreate, cmdCreate.GetTrigger())
	menuContainer.AddCommand(cmdRemoveRoleMenu, cmdRemoveRoleMenu.GetTrigger())
	menuContainer.AddCommand(cmdUpdate, cmdUpdate.GetTrigger())
	menuContainer.AddCommand(cmdResetReactions, cmdResetReactions.GetTrigger())
	menuContainer.AddCommand(cmdEditOption, cmdEditOption.GetTrigger())
	menuContainer.AddCommand(cmdFinishSetup, cmdFinishSetup.GetTrigger())
}

type ScheduledMemberRoleRemoveData struct {
	GuildID int64 `json:"guild_id"`
	GroupID int64 `json:"group_id"`
	UserID  int64 `json:"user_id"`
	RoleID  int64 `json:"role_id"`
}

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, handleReactionAddRemove, eventsystem.EventMessageReactionAdd, eventsystem.EventMessageReactionRemove)
	eventsystem.AddHandlerAsyncLastLegacy(p, handleMessageRemove, eventsystem.EventMessageDelete, eventsystem.EventMessageDeleteBulk)

	scheduledevents2.RegisterHandler("remove_member_role", ScheduledMemberRoleRemoveData{}, handleRemoveMemberRole)
	pubsub.AddHandler("role_commands_evict_menus", func(evt *pubsub.Event) {
		ClearRolemenuCache(evt.TargetGuildInt)
		recentMenusTracker.GuildReset(evt.TargetGuildInt)
	}, nil)
}

func CmdFuncRole(parsed *dcmd.Data) (interface{}, error) {
	if parsed.Args[0].Value == nil {
		return CmdFuncListCommands(parsed)
	}

	given, err := FindToggleRole(parsed.Context(), parsed.MS, parsed.Args[0].Str())
	if err != nil {
		if err == sql.ErrNoRows {
			resp, err := CmdFuncListCommands(parsed)
			if v, ok := resp.(string); ok {
				return "Não consegui encontrar o cargo, " + v, err
			}

			return resp, err
		}

		return HumanizeAssignError(parsed.GS, err)
	}

	go analytics.RecordActiveUnit(parsed.GS.ID, &Plugin{}, "cmd_used")

	if given {
		return "Te dei o cargo!", nil
	}

	return "Tirei o seu cargo!", nil
}

func HumanizeAssignError(guild *dstate.GuildState, err error) (string, error) {
	if IsRoleCommandError(err) {
		if roleError, ok := err.(*RoleError); ok {
			guild.RLock()
			defer guild.RUnlock()

			return roleError.PrettyError(guild.Guild.Roles), nil
		}

		return err.Error(), nil
	}

	if code, msg := common.DiscordError(err); code != 0 {
		switch code {
		case discordgo.ErrCodeMissingPermissions:
			return "O bot está abaixo desse cargo, fale com o administrador do servidor.", err
		case discordgo.ErrCodeMissingAccess:
			return "O bot não tem permissão suficiente para te dar esse cargo, fale com o administrador do servidor", err
		default:
			return "Ocorreu um erro ao tentar te dar o cargo: " + msg, err
		}
	}

	return "Ocorreu um erro ao tentar te dar o cargo", err

}

func CmdFuncListCommands(parsed *dcmd.Data) (interface{}, error) {
	_, grouped, ungrouped, err := GetAllRoleCommandsSorted(parsed.Context(), parsed.GS.ID)
	if err != nil {
		return "Failed retrieving role commands", err
	}

	var output strings.Builder
	output.WriteString("Here is a list of available roles:\n")

	didListCommands := false
	for group, cmds := range grouped {
		if len(cmds) < 1 {
			continue
		}
		didListCommands = true

		output.WriteString("**" + group.Name + "**\n" + StringCommands(cmds) + "\n")
	}

	if len(ungrouped) > 0 {
		didListCommands = true

		output.WriteString("**Ungrouped roles**\n" + StringCommands(ungrouped))
	}

	if !didListCommands {
		output.WriteString("No role commands (self assignable roles) set up. You can set them up in the control panel.")
	}

	return output.String(), nil
}

// StringCommands pretty formats a bunch of commands into  a string
func StringCommands(cmds []*models.RoleCommand) string {
	stringedCommands := make([]int64, 0, len(cmds))

	var output strings.Builder
	output.WriteString("```\n")

	for _, cmd := range cmds {
		if common.ContainsInt64Slice(stringedCommands, cmd.Role) {
			continue
		}

		output.WriteString(cmd.Name)
		// Check for duplicate roles
		for _, cmd2 := range cmds {
			if cmd.Role == cmd2.Role && cmd.Name != cmd2.Name {
				output.WriteString("/ " + cmd2.Name)
			}
		}
		output.WriteString("\n")

		stringedCommands = append(stringedCommands, cmd.Role)
	}

	output.WriteString("```\n")

	return output.String()
}

func handleRemoveMemberRole(evt *schEvtsModels.ScheduledEvent, data interface{}) (retry bool, err error) {
	dataCast := data.(*ScheduledMemberRoleRemoveData)
	err = common.BotSession.GuildMemberRoleRemove(dataCast.GuildID, dataCast.UserID, dataCast.RoleID)
	if err != nil {
		return scheduledevents2.CheckDiscordErrRetry(err), err
	}

	// remove the reaction
	menus, err := models.RoleMenus(
		qm.Where("role_group_id = ? AND guild_id =?", dataCast.GroupID, dataCast.GuildID),
		qm.OrderBy("message_id desc"),
		qm.Limit(10),
		qm.Load("RoleMenuOptions.RoleCommand")).AllG(context.Background())
	if err != nil {
		return false, err
	}

OUTER:
	for _, v := range menus {
		for _, opt := range v.R.RoleMenuOptions {
			if opt.R.RoleCommand.Role == dataCast.RoleID {
				// remove it
				emoji := opt.UnicodeEmoji
				if opt.EmojiID != 0 {
					emoji = "aaa:" + discordgo.StrID(opt.EmojiID)
				}

				err := common.BotSession.MessageReactionRemove(v.ChannelID, v.MessageID, emoji, dataCast.UserID)
				common.LogIgnoreError(err, "rolecommands: failed removing reaction", logrus.Fields{"guild": dataCast.GuildID, "user": dataCast.UserID, "emoji": emoji})
				continue OUTER
			}
		}
	}

	return scheduledevents2.CheckDiscordErrRetry(err), err
}

type MenuCacheKey int64

func GetRolemenuCached(ctx context.Context, gs *dstate.GuildState, messageID int64) (*models.RoleMenu, error) {
	result, err := gs.UserCacheFetch(MenuCacheKey(messageID), func() (interface{}, error) {
		menu, err := FindRolemenuFull(ctx, messageID, gs.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}

			return nil, nil
		}

		return menu, nil
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(*models.RoleMenu), nil
}

func ClearRolemenuCache(gID int64) {
	gs := bot.State.Guild(true, gID)
	if gs != nil {
		ClearRolemenuCacheGS(gs)
	}
}

func ClearRolemenuCacheGS(gs *dstate.GuildState) {
	gs.UserCacheDellAllKeysType(MenuCacheKey(0))
}
