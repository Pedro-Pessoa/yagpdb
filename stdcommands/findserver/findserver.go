package findserver

import (
	"strconv"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Pedro-Pessoa/tidbot/bot/models"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

type Candidate struct {
	ID   int64
	Name string

	UserMatch bool
	Owner     bool
	Admin     bool
	Mod       bool
}

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "findserver",
	Aliases:              []string{"findservers"},
	Description:          "Looks for a server by server name or the servers a user owns",
	HideFromHelp:         true,
	IsModCmd:             true,
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "name", Name: "name", Type: dcmd.String, Default: ""},
		{Switch: "user", Name: "user", Type: dcmd.UserID, Default: 0},
	},
	RunFunc: util.RequireBotAdmin(func(data *dcmd.Data) (interface{}, error) {
		nameToMatch := strings.ToLower(data.Switch("name").Str())
		userIDToMatch := data.Switch("user").Int64()

		if userIDToMatch == 0 && nameToMatch == "" {
			return "-name or -user not provided", nil
		}

		var whereQM qm.QueryMod
		if userIDToMatch != 0 {
			whereQM = qm.Where("owner_id = ?", userIDToMatch)
		} else {
			whereQM = qm.Where("name ILIKE ?", "%"+nameToMatch+"%")
		}

		results, err := models.JoinedGuilds(qm.Where("left_at is null"), whereQM, qm.OrderBy("id desc"), qm.Limit(250)).AllG(data.Context())
		if err != nil {
			return nil, err
		}

		var resp strings.Builder
		for _, v := range results {
			resp.WriteString("`" + strconv.FormatInt(v.ID, 10) + "`: **" + v.Name + "**\n")
		}

		resp.WriteString(strconv.Itoa(len(results)) + " results")

		return resp.String(), nil
	}),
}

func CheckGuild(gs *dstate.GuildState, nameToMatch string, userToMatch int64) *Candidate {
	if nameToMatch != "" {
		gl := strings.ToLower(gs.Guild.Name)
		if gl != nameToMatch && !strings.Contains(gl, nameToMatch) {
			return nil
		}
	}

	foundUser := false
	if userToMatch != 0 {
		for _, ms := range gs.Members {
			if ms.ID == userToMatch {
				foundUser = true
				break
			}
		}

		if !foundUser {
			return nil
		}
	}

	candidate := &Candidate{
		ID:   gs.ID,
		Name: gs.Guild.Name,
	}

	if foundUser {
		if gs.Guild.OwnerID == userToMatch {
			candidate.Owner = true
		}

		perms, _ := gs.MemberPermissions(false, 0, userToMatch)
		if perms&discordgo.PermissionAdministrator != 0 {
			candidate.Admin = true
		}

		if perms&discordgo.PermissionManageServer != 0 || perms&discordgo.PermissionKickMembers != 0 || perms&discordgo.PermissionBanMembers != 0 {
			candidate.Mod = true
		}

		candidate.UserMatch = true
	}

	return candidate
}
