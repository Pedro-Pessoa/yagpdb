package edittargetnickname

import (
	"fmt"
	"strings"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTool,
	Name:         "EditTargetNickname",
	Aliases:      []string{"etn"},
	Description:  "Edits the nickname of the specified user",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Usuário", Type: dcmd.UserID},
		{Name: "Nick", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if !bot.IsGuildWhiteListed(data.GS.ID) {
			return "Esse servidor não pode usar esse comando.", nil
		}

		if util.IsExecedByCC(data) {
			return "O comando EditTargetNickname não pode ser usado através de um CC.", nil
		}

		if ok, err := bot.AdminOrPermMS(data.CS.ID, data.MS, discordgo.PermissionManageNicknames); err != nil {
			return "Falha ao checar as permissões", err
		} else if !ok {
			return "Você precisa da permissão de gerenciar nicknames para usar esse comando!", nil
		}

		ms, err := bot.GetMember(data.GS.ID, data.Args[0].Int64())
		if err != nil {
			return "Membro não encontrado.", nil
		}

		nick := SafeArgString(data, 1)
		if strings.Compare(ms.Nick, nick) == 0 {
			return "Esse já é o apelido do usuário.", nil
		}

		err = common.BotSession.GuildMemberNickname(data.GS.ID, ms.ID, nick)
		if err != nil {
			return "", err
		}

		if nick == "" {
			return fmt.Sprintf("O apelido do usuário <@%d> foi removido.", ms.ID), nil
		}

		return fmt.Sprintf("O apelido do usuário <@%d> foi alterado para `%s`.", ms.ID, nick), nil
	},
}

func SafeArgString(data *dcmd.Data, arg int) string {
	if arg >= len(data.Args) || data.Args[arg].Value == nil {
		return ""
	}

	return data.Args[arg].Str()
}
