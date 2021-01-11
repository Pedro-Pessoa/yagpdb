package moderation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/analytics"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/paginatedmessages"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/scheduledevents2"
)

func MBaseCmd(cmdData *dcmd.Data, targetID int64) (config *Config, targetUser *discordgo.User, err error) {
	config, err = GetConfig(cmdData.GS.ID)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "GetConfig")
	}

	if targetID != 0 {
		targetMember, _ := bot.GetMember(cmdData.GS.ID, targetID)
		if targetMember != nil {
			gs := cmdData.GS

			gs.RLock()
			above := bot.IsMemberAbove(gs, cmdData.MS, targetMember)
			gs.RUnlock()

			if !above {
				return config, targetMember.DGoUser(), commands.NewUserError("Você não pode usar comandos de moderação em usuários acima de você.")
			}

			return config, targetMember.DGoUser(), nil
		}
	}

	return config, &discordgo.User{
		Username:      "unknown",
		Discriminator: "????",
		ID:            targetID,
	}, nil

}

func MBaseCmdSecond(cmdData *dcmd.Data, reason string, reasonArgOptional bool, neededPerm int, additionalPermRoles []int64, enabled bool) (oreason string, err error) {
	cmdName := cmdData.Cmd.Trigger.Names[0]
	oreason = reason
	if !enabled {
		return oreason, commands.NewUserErrorf("O comando **%s** está desativado.", cmdName)
	}

	if strings.TrimSpace(reason) == "" {
		if !reasonArgOptional {
			return oreason, commands.NewUserError("Você precisa dizer o motivo pela qual está usando esse comando.")
		}

		oreason = "(Sem motivo especificado)"
	}

	member := cmdData.MS

	// check permissions or role setup for this command
	permsMet := false
	if len(additionalPermRoles) > 0 {
		// Check if the user has one of the required roles
		for _, r := range member.Roles {
			if common.ContainsInt64Slice(additionalPermRoles, r) {
				permsMet = true
				break
			}
		}
	}

	if !permsMet && neededPerm != 0 {
		// Fallback to legacy permissions
		hasPerms, err := bot.AdminOrPermMS(cmdData.CS.ID, member, neededPerm)
		if err != nil || !hasPerms {
			return oreason, commands.NewUserErrorf("O comando **%s** exige permissão de **%s** nesse canal.", cmdName, common.StringPerms[neededPerm])
		}
	}

	go analytics.RecordActiveUnit(cmdData.GS.ID, &Plugin{}, "executed_cmd_"+cmdName)

	return oreason, nil
}

func SafeArgString(data *dcmd.Data, arg int) string {
	if arg >= len(data.Args) || data.Args[arg].Value == nil {
		return ""
	}

	return data.Args[arg].Str()
}

func GenericCmdResp(action ModlogAction, target *discordgo.User, duration time.Duration, zeroDurPermanent bool, noDur bool) string {
	durStr := " sem prazo para acabar!"
	if duration > 0 || !zeroDurPermanent {
		durStr = " por `" + common.HumanizeDuration(common.DurationPrecisionMinutes, duration) + "`"
	}

	if noDur {
		durStr = ""
	}

	userStr := target.Username + "#" + target.Discriminator
	if target.Discriminator == "????" {
		userStr = strconv.FormatInt(target.ID, 10)
	}

	return fmt.Sprintf("%s %s `%s`%s", action.Emoji, action.Prefix, userStr, durStr)
}

var ModerationCommands = []*commands.YAGCommand{
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Ban",
		Aliases:       []string{"banid", "banir"},
		Description:   "Bane um memebro, especifique uma duração com -d e especifique o número de dias de mensagens para deletar com -ddays (de 0 a 7)",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "d", Default: time.Duration(0), Name: "Duration", Type: &commands.DurationArg{}},
			{Switch: "ddays", Default: 1, Name: "Days", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			reason := SafeArgString(parsed, 1)
			reason, err = MBaseCmdSecond(parsed, reason, config.BanReasonOptional, discordgo.PermissionBanMembers, config.BanCmdRoles, config.BanEnabled)
			if err != nil {
				return nil, err
			}

			ddays := int(config.DefaultBanDeleteDays.Int64)
			if parsed.Switches["ddays"].Value != nil {
				ddays = parsed.Switches["ddays"].Int()
			}

			err = BanUserWithDuration(config, parsed.GS.ID, parsed.CS, parsed.Msg, parsed.Msg.Author, reason, target, parsed.Switches["d"].Value.(time.Duration), ddays)
			if err != nil {
				return nil, err
			}

			return GenericCmdResp(MABanned, target, parsed.Switch("d").Value.(time.Duration), true, false), nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Unban",
		Aliases:       []string{"unbanid", "desbanir"},
		Description:   "Desbane um usuário.",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},

		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0) //in most situations, the target will not be a part of server, hence no point in doing unnecessary api calls(i.e. bot.GetMember)
			if err != nil {
				return nil, err
			}

			reason := SafeArgString(parsed, 1)
			reason, err = MBaseCmdSecond(parsed, reason, config.BanReasonOptional, discordgo.PermissionBanMembers, config.BanCmdRoles, config.BanEnabled)
			if err != nil {
				return nil, err
			}
			targetID := parsed.Args[0].Int64()
			target := &discordgo.User{
				Username:      "unknown",
				Discriminator: "????",
				ID:            targetID,
			}

			targetMem := parsed.GS.MemberCopy(true, targetID)
			if targetMem != nil {
				return "Esse usuário não está banido!", nil
			}

			isNotBanned, err := UnbanUser(config, parsed.GS.ID, parsed.Msg.Author, reason, target)

			if err != nil {
				return nil, err
			}
			if isNotBanned {
				return "Esse usuário não está banido!", nil
			}

			return GenericCmdResp(MAUnbanned, target, 0, true, true), nil
		},
	},
	{
		CmdCategory:        commands.CategoryModeration,
		Name:               "LockDown",
		Aliases:            []string{"ld", "lock", "trancar"},
		Description:        "Bloqueia o servidor ou algum cargo específico.",
		LongDescription:    "Requer permissão de gerenciamento de cargos. Esse comando vai retirar a permissão do \"everyone\" de enviar mensagens.\nVocê pode escolher um cargo para ser bloqueado usando o nome ou o ID dele.\n\nVocê também pode usar flags para retirar mais permissões:\n**-reaction** -> Retira a permissão de adicionar reações\n**-voicespeak** -> Retira a permissão de falar\n**-voiceconnect** -> Retira a permissão de se conectar a uma canal de voz\n**-all* -> Retira todas as permissões anteriores\n**-force** -> As permissões originais do cargo são sobrepostas durante o unlock",
		GuildScopeCooldown: 10,
		RequiredArgs:       0,
		Arguments: []*dcmd.ArgDef{
			{Name: "Cargo", Help: "Cargo opcional", Type: dcmd.String},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "reaction", Name: "Reações"},
			{Switch: "voicespeak", Name: "Permissão de fala"},
			{Switch: "voiceconnect", Name: "Permissão de conexão"},
			{Switch: "all", Name: "Todas as Flags"},
			{Switch: "force", Name: "Força sobreposição de permissões", Default: false},
			{Switch: "d", Name: "Duração", Type: &commands.DurationArg{}},
		},
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(data, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(data, "", true, discordgo.PermissionManageRoles, config.LockdownCmdRoles, config.LockdownCmdEnabled)
			if err != nil {
				return nil, err
			}

			totalPerms := discordgo.PermissionSendMessages

			if data.Switches["all"].Value != nil && data.Switches["all"].Value.(bool) {
				totalPerms = totalPerms | discordgo.PermissionAddReactions | discordgo.PermissionVoiceSpeak | discordgo.PermissionVoiceConnect
			} else {
				if data.Switches["reaction"].Value != nil && data.Switches["reaction"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionAddReactions
				}

				if data.Switches["voicespeak"].Value != nil && data.Switches["voicespeak"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionVoiceSpeak
				}

				if data.Switches["voiceconnect"].Value != nil && data.Switches["voiceconnect"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionVoiceConnect
				}
			}

			dur := time.Duration(config.DefaultLockdownDuration.Int64) * time.Minute
			if d := data.Switches["d"].Value; d != nil {
				dur = d.(time.Duration)
			}

			out, err := LockUnlockRole(config, true, data.GS, data.CS, data.MS, data.Msg.Author, "Moderação", data.Args[0].Str(), data.Switches["force"].Value.(bool), totalPerms, dur)
			if err != nil {
				return nil, err
			}

			return out, nil
		},
	},
	{
		CmdCategory:        commands.CategoryModeration,
		Name:               "UnLock",
		Aliases:            []string{"ul", "destrancar"},
		Description:        "Unlocks the server or a specific role.",
		LongDescription:    "Requer permissão de gerenciamento de cargos.\nEsse comando vai adicionar a permissão do \"everyone\" de enviar mensagens.\nVocê pode escolher um cargo para ser bloqueado usando o nome ou o ID dele.\n\nVocê também pode usar flags para adicionar mais permissões:\n**-reaction** -> Adiciona a permissão de adicionar reações\n**-voicespeak** -> Adiciona a permissão de falar\n**-voiceconnect** -> Adiciona a permissão de se conectar a uma canal de voz\n**-all* -> Adiciona todas as permissões anteriores\n**-force** -> As permissões originais do cargo são sobrepostas.",
		GuildScopeCooldown: 10,
		RequiredArgs:       0,
		Arguments: []*dcmd.ArgDef{
			{Name: "Cargo", Help: "Optional role", Type: dcmd.String},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "reaction", Name: "Reações"},
			{Switch: "voicespeak", Name: "Permissão de fala"},
			{Switch: "voiceconnect", Name: "Permissão de conexão"},
			{Switch: "all", Name: "Todas as Flags"},
			{Switch: "force", Name: "Força sobreposição de permissões", Default: false},
		},
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(data, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(data, "", true, discordgo.PermissionManageRoles, config.LockdownCmdRoles, config.LockdownCmdEnabled)
			if err != nil {
				return nil, err
			}

			totalPerms := discordgo.PermissionSendMessages

			if data.Switches["all"].Value != nil && data.Switches["all"].Value.(bool) {
				totalPerms = totalPerms | discordgo.PermissionAddReactions | discordgo.PermissionVoiceSpeak | discordgo.PermissionVoiceConnect
			} else {
				if data.Switches["reaction"].Value != nil && data.Switches["reaction"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionAddReactions
				}

				if data.Switches["voicespeak"].Value != nil && data.Switches["voicespeak"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionVoiceSpeak
				}

				if data.Switches["voiceconnect"].Value != nil && data.Switches["voiceconnect"].Value.(bool) {
					totalPerms = totalPerms | discordgo.PermissionVoiceConnect
				}
			}

			out, err := LockUnlockRole(config, false, data.GS, data.CS, data.MS, data.Msg.Author, "Moderação", data.Args[0].Str(), data.Switches["force"].Value.(bool), totalPerms, time.Duration(0))
			if err != nil {
				return nil, err
			}

			return out, nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Kick",
		Aliases:       []string{"kikar", "expulsar"},
		Description:   "Expulsa um membro",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			reason := SafeArgString(parsed, 1)
			reason, err = MBaseCmdSecond(parsed, reason, config.KickReasonOptional, discordgo.PermissionKickMembers, config.KickCmdRoles, config.KickEnabled)
			if err != nil {
				return nil, err
			}

			err = KickUser(config, parsed.GS.ID, parsed.CS, parsed.Msg, parsed.Msg.Author, reason, target)
			if err != nil {
				return nil, err
			}

			return GenericCmdResp(MAKick, target, 0, true, true), nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Mute",
		Aliases:       []string{"silenciar", "mutar"},
		Description:   "Silencia um membro",
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Duração", Type: &commands.DurationArg{}},
			{Name: "Motivo", Type: dcmd.String},
		},
		ArgumentCombos: [][]int{{0, 1, 2}, {0, 2, 1}, {0, 1}, {0, 2}, {0}},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			if config.MuteRole == "" {
				return "O cargo de silenciado não foi configurado, por favor configure no painel de controle.", nil
			}

			reason := parsed.Args[2].Str()
			reason, err = MBaseCmdSecond(parsed, reason, config.MuteReasonOptional, discordgo.PermissionKickMembers, config.MuteCmdRoles, config.MuteEnabled)
			if err != nil {
				return nil, err
			}

			d := time.Duration(config.DefaultMuteDuration.Int64) * time.Minute
			if parsed.Args[1].Value != nil {
				d = parsed.Args[1].Value.(time.Duration)
			}

			if d > 0 && d < time.Minute {
				d = time.Minute
			}

			logger.Info(d.Seconds())

			member, err := bot.GetMember(parsed.GS.ID, target.ID)
			if err != nil || member == nil {
				return "Membro não encontrado", err
			}

			err = MuteUnmuteUser(config, true, parsed.GS.ID, parsed.CS, parsed.Msg, parsed.Msg.Author, reason, member, int(d.Minutes()))
			if err != nil {
				return nil, err
			}

			return GenericCmdResp(MAMute, target, d, true, false), nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Unmute",
		Aliases:       []string{"desmutar", "desilenciar"},
		Description:   "Dessilencia um membro",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			if config.MuteRole == "" {
				return "O cargo de silenciado não foi configurado, por favor configure no painel de controle.", nil
			}

			reason := parsed.Args[1].Str()
			reason, err = MBaseCmdSecond(parsed, reason, config.UnmuteReasonOptional, discordgo.PermissionKickMembers, config.MuteCmdRoles, config.MuteEnabled)
			if err != nil {
				return nil, err
			}

			member, err := bot.GetMember(parsed.GS.ID, target.ID)
			if err != nil || member == nil {
				return "Membro não encontrado.", err
			}

			err = MuteUnmuteUser(config, false, parsed.GS.ID, parsed.CS, parsed.Msg, parsed.Msg.Author, reason, member, 0)
			if err != nil {
				return nil, err
			}

			return GenericCmdResp(MAUnmute, target, 0, false, true), nil
		},
	},
	{
		CustomEnabled: true,
		Cooldown:      5,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Report",
		Aliases:       []string{"reportar"},
		Description:   "Reporta um membro para o staff do servidor.",
		RequiredArgs:  2,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, 0, nil, config.ReportEnabled)
			if err != nil {
				return nil, err
			}

			target := parsed.Args[0].Int64()

			if target == parsed.Msg.Author.ID {
				return "Você não pode se reportar, bobão.", nil
			}

			logLink := CreateLogs(parsed.GS.ID, parsed.CS.ID, parsed.Msg.Author)

			channelID := config.IntReportChannel()
			if channelID == 0 {
				return "O canel de reports não foi configurado.", nil
			}

			reportBody := fmt.Sprintf("<@%d> Reportou <@%d> in <#%d> For `%s`\nÚltimas 100 mensagens do canal: <%s>", parsed.Msg.Author.ID, target, parsed.Msg.ChannelID, parsed.Args[1].Str(), logLink)

			_, err = common.BotSession.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
				Content: reportBody,
				AllowedMentions: discordgo.AllowedMentions{
					Users: []int64{parsed.Msg.Author.ID, target},
				},
			})

			if err != nil {
				return nil, err
			}

			// don't bother sending confirmation if it's in the same channel
			if channelID != parsed.Msg.ChannelID {
				return "Usuário reportado para as autoridades!", nil
			}

			return nil, nil
		},
	},
	{
		CustomEnabled:   true,
		CmdCategory:     commands.CategoryModeration,
		Name:            "Clean",
		Description:     "Deleta as últimas mensagens do chat pela quantidade especificada (no máximo 100). Opcionalmente filtrado por usuário.",
		LongDescription: "Você também pode usar essas flags:\n**-r** Especifique um regex\n**-ma** Duração máxima\n**-minage** Duração mínima\nOBS: Somente as últimas 1000 mensagens serão verificadas.\n**-i** Faz o regex ignorar capitalização\n**-nopin** Não deleta mensagens fixadas\n**-to** Para a execução do clean quando chegar nessa mensagem.",
		Aliases:         []string{"clear", "cl", "limpar"},
		RequiredArgs:    1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Número", Type: &dcmd.IntArg{Min: 1, Max: 100}},
			{Name: "Usuário", Type: dcmd.UserID, Default: 0},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "r", Name: "Regex", Type: dcmd.String},
			{Switch: "ma", Default: time.Duration(0), Name: "Max age", Type: &commands.DurationArg{}},
			{Switch: "minage", Default: time.Duration(0), Name: "Min age", Type: &commands.DurationArg{}},
			{Switch: "i", Name: "Regex case insensitive"},
			{Switch: "nopin", Name: "Ignore pinned messages"},
			{Switch: "to", Name: "Stop at this msg ID", Type: dcmd.Int},
		},
		ArgumentCombos: [][]int{{0}, {0, 1}, {1, 0}},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, nil, config.CleanEnabled)
			if err != nil {
				return nil, err
			}

			userFilter := parsed.Args[1].Int64()

			num := parsed.Args[0].Int()
			if (userFilter == 0 || userFilter == parsed.Msg.Author.ID) && parsed.Source != 0 {
				num++ // Automatically include our own message if not triggeded by exec/execAdmin
			}

			if num > 100 {
				num = 100
			}

			if num < 1 {
				if num < 0 {
					return errors.New("O bot não está se sentindo bem"), nil
				}
				return errors.New("Não deu pra deletar nada!"), nil
			}

			filtered := false

			// Check if we should regex match this
			re := ""
			if parsed.Switches["r"].Value != nil {
				filtered = true
				re = parsed.Switches["r"].Str()

				// Add the case insensitive flag if needed
				if parsed.Switches["i"].Value != nil && parsed.Switches["i"].Value.(bool) {
					if !strings.HasPrefix(re, "(?i)") {
						re = "(?i)" + re
					}
				}
			}

			// Check if we have a max age
			ma := parsed.Switches["ma"].Value.(time.Duration)
			if ma != 0 {
				filtered = true
			}

			// Check if we have a min age
			minAge := parsed.Switches["minage"].Value.(time.Duration)
			if minAge != 0 {
				filtered = true
			}

			// Break if it gets to this msg ID
			toID := int64(0)
			if parsed.Switches["to"].Value != nil {
				filtered = true
				toID = parsed.Switches["to"].Int64()
			}

			// Check if we should ignore pinned messages
			pe := false
			if parsed.Switches["nopin"].Value != nil && parsed.Switches["nopin"].Value.(bool) {
				pe = true
				filtered = true
			}

			limitFetch := num
			if userFilter != 0 || filtered {
				limitFetch = num * 50 // Maybe just change to full fetch?
			}

			if limitFetch > 1000 {
				limitFetch = 1000
			}

			// Wait a second so the client dosen't gltich out
			time.Sleep(time.Second)

			numDeleted, err := AdvancedDeleteMessages(parsed.Msg.ChannelID, userFilter, re, toID, ma, minAge, pe, num, limitFetch)

			return dcmd.NewTemporaryResponse(time.Second*5, fmt.Sprintf("Apaguei %d mensagens! :')", numDeleted), true), err
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Reason",
		Aliases:       []string{"motivo"},
		Description:   "Adiciona/Edita uma razão no modlog",
		RequiredArgs:  2,
		Arguments: []*dcmd.ArgDef{
			{Name: "ID da mensagem", Type: dcmd.Int},
			{Name: "Motivo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionKickMembers, nil, true)
			if err != nil {
				return nil, err
			}

			if config.ActionChannel == "" {
				return "O canal de modlog não foi definido no painel de controle.", nil
			}

			msg, err := common.BotSession.ChannelMessage(config.IntActionChannel(), parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			if msg.Author.ID != common.BotUser.ID {
				return "Não fui eu quem fiz essa mensagem.", nil
			}

			if len(msg.Embeds) < 1 {
				return "Essa mensagem é muito antiga ou você está me zoando ;)", nil
			}

			embed := msg.Embeds[0]
			updateEmbedReason(parsed.Msg.Author, parsed.Args[1].Str(), embed)
			_, err = common.BotSession.ChannelMessageEditEmbed(config.IntActionChannel(), msg.ID, embed)
			if err != nil {
				return nil, err
			}

			return "👌", nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Warn",
		Aliases:       []string{"notificar", "avisar"},
		Description:   "Notifica um usuário. Use `-warnings` para ver as notificações.",
		RequiredArgs:  2,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Motivo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}
			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, config.WarnCommandsEnabled)
			if err != nil {
				return nil, err
			}

			err = WarnUser(config, parsed.GS.ID, parsed.CS, parsed.Msg, parsed.Msg.Author, target, parsed.Args[1].Str())
			if err != nil {
				return nil, err
			}

			return GenericCmdResp(MAWarned, target, 0, false, true), nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "Warnings",
		Aliases:       []string{"Avisos", "Warns"},
		Description:   "Lista as notificações de um usuário",
		RequiredArgs:  0,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID, Default: 0},
			{Name: "Página", Type: &dcmd.IntArg{Max: 10000}, Default: 0},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "id", Name: "ID da notificação", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			var err error
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, true)
			if err != nil {
				return nil, err
			}

			if parsed.Switches["id"].Value != nil {
				var warn []*WarningModel
				err = common.GORM.Where("guild_id = ? AND id = ?", parsed.GS.ID, parsed.Switches["id"].Int()).First(&warn).Error
				if err != nil && err != gorm.ErrRecordNotFound {
					return nil, err
				}
				if len(warn) == 0 {
					return fmt.Sprintf("O aviso com id : `%d` não existe", parsed.Switches["id"].Int()), nil
				}

				return &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("Aviso#%d - Usuário : %s", warn[0].ID, warn[0].UserID),
					Description: fmt.Sprintf("`%20s` - **Motivo** : %s", warn[0].CreatedAt.UTC().Format(time.RFC822), warn[0].Message),
					Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("By: %s (%13s)", warn[0].AuthorUsernameDiscrim, warn[0].AuthorID)},
				}, nil
			}

			page := parsed.Args[1].Int()
			if page < 1 {
				page = 1
			}

			if parsed.Context().Value(paginatedmessages.CtxKeyNoPagination) != nil {
				return PaginateWarnings(parsed)(nil, page)
			}

			_, err = paginatedmessages.CreatePaginatedMessage(parsed.GS.ID, parsed.CS.ID, page, 0, PaginateWarnings(parsed))
			return nil, err
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "EditWarning",
		Aliases:       []string{"EditarWarning", "EditarNotificacao"},
		Description:   "Edita uma notificação. O ID é o primeiro número de cada notificação no comando `warnings`.",
		RequiredArgs:  2,
		Arguments: []*dcmd.ArgDef{
			{Name: "ID", Type: dcmd.Int},
			{Name: "NovaMensagem", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, config.WarnCommandsEnabled)
			if err != nil {
				return nil, err
			}

			rows := common.GORM.Model(WarningModel{}).Where("guild_id = ? AND id = ?", parsed.GS.ID, parsed.Args[0].Int()).Update(
				"message", fmt.Sprintf("%s (modificado por %s#%s (%d))", parsed.Args[1].Str(), parsed.Msg.Author.Username, parsed.Msg.Author.Discriminator, parsed.Msg.Author.ID)).RowsAffected

			if rows < 1 {
				return "Falha ao atualizar, provavelmente não consegui encontrar a notificação.", nil
			}

			return "👌", nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "DelWarning",
		Aliases:       []string{"dw"},
		Description:   "Deleta uma notificação. O ID é o primeiro número de cada notificação no comando `warnings`.",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "ID", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, config.WarnCommandsEnabled)
			if err != nil {
				return nil, err
			}

			rows := common.GORM.Where("guild_id = ? AND id = ?", parsed.GS.ID, parsed.Args[0].Int()).Delete(WarningModel{}).RowsAffected
			if rows < 1 {
				return "Falha ao deletar, provavelmente não consegui encontrar a notificação.", nil
			}

			return "👌", nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "ClearWarnings",
		Aliases:       []string{"clw"},
		Description:   "Apaga todas as notificações do usuário.",
		RequiredArgs:  1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {

			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, config.WarnCommandsEnabled)
			if err != nil {
				return nil, err
			}

			userID := parsed.Args[0].Int64()

			rows := common.GORM.Where("guild_id = ? AND user_id = ?", parsed.GS.ID, userID).Delete(WarningModel{}).RowsAffected
			return fmt.Sprintf("Deletei %d notificações.", rows), nil
		},
	},
	{
		CmdCategory:  commands.CategoryModeration,
		Name:         "Slowmode",
		Aliases:      []string{"sm"},
		Description:  "Muda o slowmode do canal durante uma duração opcional.",
		RequiredArgs: 1,
		Cooldown:     10,
		Arguments: []*dcmd.ArgDef{
			{Name: "RateLimit", Type: dcmd.Int},
			{Name: "Duração", Type: &commands.DurationArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, nil, config.SlowmodeCommandEnabled)
			if err != nil {
				return nil, err
			}

			duration := time.Duration(0)
			if parsed.Args[1].Value != nil {
				duration = parsed.Args[1].Value.(time.Duration)
			}

			if duration > 0 && duration < time.Minute {
				duration = time.Minute
			}

			RL := parsed.Args[0].Int()
			if RL > 21600 {
				RL = 21600
			}

			out, err := SlowModeFunc(config, parsed.GS.ID, parsed.CS, parsed.Msg.Author, int(duration.Minutes()), RL)
			if err != nil {
				return nil, err
			}

			return out, nil
		},
	},
	{
		CmdCategory: commands.CategoryModeration,
		Name:        "TopWarnings",
		Aliases:     []string{"topwarns", "topnotificacoes"},
		Description: "Mostra uma lista de notificaçoes do server.",
		Arguments: []*dcmd.ArgDef{
			{Name: "Página", Type: dcmd.Int, Default: 0},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "ID", Name: "List userIDs"},
		},
		RunFunc: paginatedmessages.PaginatedCommand(0, func(parsed *dcmd.Data, p *paginatedmessages.PaginatedMessage, page int) (*discordgo.MessageEmbed, error) {
			showUserIDs := false
			config, _, err := MBaseCmd(parsed, 0)
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageMessages, config.WarnCmdRoles, true)
			if err != nil {
				return nil, err
			}

			if parsed.Switches["id"].Value != nil && parsed.Switches["id"].Value.(bool) {
				showUserIDs = true
			}

			offset := (page - 1) * 15
			entries, err := TopWarns(parsed.GS.ID, offset, 15)
			if err != nil {
				return nil, err
			}

			if len(entries) < 1 && p != nil && p.LastResponse != nil { //Don't send No Results error on first execution.
				return nil, paginatedmessages.ErrNoResults
			}

			embed := &discordgo.MessageEmbed{
				Title: "Lista ranqueada de notificações",
			}

			out := "```\n# - Warns - User\n"
			for _, v := range entries {
				if !showUserIDs {
					user := v.Username
					if user == "" {
						user = "unknown ID:" + strconv.FormatInt(v.UserID, 10)
					}
					out += fmt.Sprintf("#%02d: %4d - %s\n", v.Rank, v.WarnCount, user)
				} else {
					out += fmt.Sprintf("#%02d: %4d - %d\n", v.Rank, v.WarnCount, v.UserID)
				}
			}
			out += "```\n"

			embed.Description = out

			return embed, nil

		}),
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "GiveRole",
		Aliases:       []string{"grole", "arole", "addrole"},
		Description:   "Dá um cargo para o usuário especificado, com uma duração opcional.",

		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Cargo", Type: dcmd.String},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "d", Default: time.Duration(0), Name: "Duração", Type: &commands.DurationArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageRoles, config.GiveRoleCmdRoles, config.GiveRoleCmdEnabled)
			if err != nil {
				return nil, err
			}

			member, err := bot.GetMember(parsed.GS.ID, target.ID)
			if err != nil || member == nil {
				return "Membro não encontrado.", err
			}

			role := FindRole(parsed.GS, parsed.Args[1].Str())
			if role == nil {
				return "Não encontrei o cargo especificado.", nil
			}

			parsed.GS.RLock()
			if !bot.IsMemberAboveRole(parsed.GS, parsed.MS, role) {
				parsed.GS.RUnlock()
				return "Não pode dar cargos acima dos seus.", nil
			}
			parsed.GS.RUnlock()

			dur := parsed.Switches["d"].Value.(time.Duration)

			// no point if the user has the role and is not updating the expiracy
			if common.ContainsInt64Slice(member.Roles, role.ID) && dur <= 0 {
				return "Esse usuário já tem esse cargo.", nil
			}

			err = common.AddRoleDS(member, role.ID)
			if err != nil {
				return nil, err
			}

			// schedule the expirey
			if dur > 0 {
				err := scheduledevents2.ScheduleRemoveRole(parsed.Context(), parsed.GS.ID, target.ID, role.ID, time.Now().Add(dur))
				if err != nil {
					return nil, err
				}
			}

			// cancel the event to add the role
			_ = scheduledevents2.CancelAddRole(parsed.Context(), parsed.GS.ID, parsed.Msg.Author.ID, role.ID)

			action := MAGiveRole
			action.Prefix = "Cargo " + role.Name + " adicionado a(o) "
			if config.GiveRoleCmdModlog && config.IntActionChannel() != 0 {
				if dur > 0 {
					action.Footer = "Duração: " + common.HumanizeDuration(common.DurationPrecisionMinutes, dur)
				}
				_ = CreateModlogEmbed(config, parsed.Msg.Author, action, target, "", "")
			}

			return GenericCmdResp(action, target, dur, true, dur <= 0), nil
		},
	},
	{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryModeration,
		Name:          "RemoveRole",
		Aliases:       []string{"rrole", "takerole", "trole", "tirarcargo"},
		Description:   "Retira o cargo especificado do usuário especificado.",

		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			{Name: "Usuário", Type: dcmd.UserID},
			{Name: "Cargo", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			config, target, err := MBaseCmd(parsed, parsed.Args[0].Int64())
			if err != nil {
				return nil, err
			}

			_, err = MBaseCmdSecond(parsed, "", true, discordgo.PermissionManageRoles, config.GiveRoleCmdRoles, config.GiveRoleCmdEnabled)
			if err != nil {
				return nil, err
			}

			member, err := bot.GetMember(parsed.GS.ID, target.ID)
			if err != nil || member == nil {
				return "Membro não encontrado", err
			}

			role := FindRole(parsed.GS, parsed.Args[1].Str())
			if role == nil {
				return "Cargo não encontrado", nil
			}

			parsed.GS.RLock()
			if !bot.IsMemberAboveRole(parsed.GS, parsed.MS, role) {
				parsed.GS.RUnlock()
				return "Você não pode remover cargos mais altos que o seu.", nil
			}
			parsed.GS.RUnlock()

			err = common.RemoveRoleDS(member, role.ID)
			if err != nil {
				return nil, err
			}

			// cancel the event to remove the role
			_ = scheduledevents2.CancelRemoveRole(parsed.Context(), parsed.GS.ID, parsed.Msg.Author.ID, role.ID)

			action := MARemoveRole
			action.Prefix = "Cargo " + role.Name + " removido de "
			if config.GiveRoleCmdModlog && config.IntActionChannel() != 0 {
				_ = CreateModlogEmbed(config, parsed.Msg.Author, action, target, "", "")
			}

			return GenericCmdResp(action, target, 0, true, true), nil
		},
	},
}

func AdvancedDeleteMessages(channelID int64, filterUser int64, regex string, toID int64, maxAge time.Duration, minAge time.Duration, pinFilterEnable bool, deleteNum, fetchNum int) (int, error) {
	var compiledRegex *regexp.Regexp
	if regex != "" {
		// Start by compiling the regex
		var err error
		compiledRegex, err = regexp.Compile(regex)
		if err != nil {
			return 0, err
		}
	}

	var pinnedMessages map[int64]struct{}
	if pinFilterEnable {
		//Fetch pinned messages from channel and make a map with ids as keys which will make it easy to verify if a message with a given ID is pinned message
		messageSlice, err := common.BotSession.ChannelMessagesPinned(channelID)
		if err != nil {
			return 0, err
		}
		pinnedMessages = make(map[int64]struct{}, len(messageSlice))
		for _, msg := range messageSlice {
			pinnedMessages[msg.ID] = struct{}{} //empty struct works because we are not really interested in value
		}
	}

	msgs, err := bot.GetMessages(channelID, fetchNum, false)
	if err != nil {
		return 0, err
	}

	toDelete := make([]int64, 0)
	now := time.Now()
	for i := len(msgs) - 1; i >= 0; i-- {
		if filterUser != 0 && msgs[i].Author.ID != filterUser {
			continue
		}

		// Can only bulk delete messages up to 2 weeks (but add 1 minute buffer account for time sync issues and other smallies)
		if now.Sub(msgs[i].ParsedCreated) > (time.Hour*24*14)-time.Minute {
			continue
		}

		// Check regex
		if compiledRegex != nil {
			if !compiledRegex.MatchString(msgs[i].Content) {
				continue
			}
		}

		// Check max age
		if maxAge != 0 && now.Sub(msgs[i].ParsedCreated) > maxAge {
			continue
		}

		// Check min age
		if minAge != 0 && now.Sub(msgs[i].ParsedCreated) < minAge {
			continue
		}

		// Check if pinned message to ignore
		if pinFilterEnable {
			if _, found := pinnedMessages[msgs[i].ID]; found {
				continue
			}
		}

		// Continue only if current msg ID is < toID
		if toID > msgs[i].ID {
			break
		}

		toDelete = append(toDelete, msgs[i].ID)
		//log.Println("Deleting", msgs[i].ContentWithMentionsReplaced())
		if len(toDelete) >= deleteNum || len(toDelete) >= 100 {
			break
		}
	}

	if len(toDelete) < 1 {
		return 0, nil
	}

	if len(toDelete) < 1 {
		return 0, nil
	} else if len(toDelete) == 1 {
		err = common.BotSession.ChannelMessageDelete(channelID, toDelete[0])
	} else {
		err = common.BotSession.ChannelMessagesBulkDelete(channelID, toDelete)
	}

	return len(toDelete), err
}

func FindRole(gs *dstate.GuildState, roleS string) *discordgo.Role {
	parsedNumber, parseErr := strconv.ParseInt(roleS, 10, 64)
	var name string
	var id int64
	var err error

	if parseErr != nil { // it's a mention or a name
		if strings.HasPrefix(roleS, "<@&") && strings.HasSuffix(roleS, ">") && len(roleS) > 4 {
			id, err = strconv.ParseInt(roleS[3:len(roleS)-1], 10, 64) // if no error, it's the id
			if err != nil {                                           // it's the name
				name = roleS
			}
		}
	} else {
		id = parsedNumber
	}

	gs.RLock()
	defer gs.RUnlock()
	if name != "" { // it was a name
		for _, v := range gs.Guild.Roles {
			if strings.EqualFold(strings.TrimSpace(v.Name), name) {
				return v
			}
		}
	} else { // was a number, try looking by id
		r := gs.RoleCopy(false, id)
		if r != nil {
			return r
		}
	}

	return nil // couldn't find the role :(
}

func PaginateWarnings(parsed *dcmd.Data) func(p *paginatedmessages.PaginatedMessage, page int) (*discordgo.MessageEmbed, error) {
	return func(p *paginatedmessages.PaginatedMessage, page int) (*discordgo.MessageEmbed, error) {
		var err error
		skip := (page - 1) * 6
		userID := parsed.Args[0].Int64()
		limit := 6

		var result []*WarningModel
		var count int
		err = common.GORM.Table("moderation_warnings").Where("user_id = ? AND guild_id = ?", userID, parsed.GS.ID).Count(&count).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		err = common.GORM.Where("user_id = ? AND guild_id = ?", userID, parsed.GS.ID).Order("id desc").Offset(skip).Limit(limit).Find(&result).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if len(result) < 1 && p != nil && p.LastResponse != nil { //Dont send No Results error on first execution
			return nil, paginatedmessages.ErrNoResults
		}

		desc := fmt.Sprintf("**Total :** `%d`", count)
		var fields []*discordgo.MessageEmbedField
		currentField := &discordgo.MessageEmbedField{
			Name:  "⠀", //Use braille blank character for seamless transition between feilds
			Value: "",
		}
		fields = append(fields, currentField)
		if len(result) > 0 {
			for _, entry := range result {
				entry_formatted := fmt.Sprintf("#%d: `%20s` - Por: **%s** (%13s) \n **Motivo:** %s", entry.ID, entry.CreatedAt.UTC().Format(time.RFC822), entry.AuthorUsernameDiscrim, entry.AuthorID, entry.Message)
				if len([]rune(entry_formatted)) > 900 {
					entry_formatted = common.CutStringShort(entry_formatted, 900)
				}
				entry_formatted += "\n"
				if entry.LogsLink != "" {
					entry_formatted += fmt.Sprintf("> logs: [`link`](%s)\n", entry.LogsLink)
				}

				if len([]rune(currentField.Value+entry_formatted)) > 1023 {
					currentField = &discordgo.MessageEmbedField{
						Name:  "⠀",
						Value: entry_formatted + "\n",
					}
					fields = append(fields, currentField)
				} else {
					currentField.Value += entry_formatted + "\n"
				}
			}

		} else {
			currentField.Value = "Sem notificações"
		}

		return &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Warnings - User : %d", userID),
			Description: desc,
			Fields:      fields,
		}, nil
	}
}
