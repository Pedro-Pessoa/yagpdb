package tibiacmds

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/stdcommands/util"
	"github.com/jonas747/yagpdb/tibia"
)

var TibiaSetWorld = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "TibiaSetWorld",
	Aliases:      []string{"tsw", "mundo", "world"},
	Description:  "Determina qual vai ser o mundo deste servidor! **IRREVERSÍVEL**.",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Mundo", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.SetServerWorld(data.Args[0].Str(), data.GS.ID, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaSetGuild = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "TibiaSetGuild",
	Aliases:      []string{"tsg", "guild"},
	Description:  "Determina qual vai ser a guild deste servidor! **IRREVERSÍVEL**.",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome da Guild", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.SetServerGuild(data.Args[0].Str(), data.GS.ID, false, data.GS.Guild.MemberCount)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaSetDeathChannel = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaSetDeathChannel",
	Aliases:     []string{"tsdc", "deathchannel", "dc"},
	Description: "O canal onde esse comando for usado será utilizado para enviar avisos de morte.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.SetServerDeathChannel(data.GS.ID, data.CS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaSetUpdatesChannel = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaSetUpdatesChannel",
	Aliases:     []string{"tsuc", "updateshannel", "updatehannel", "uc"},
	Description: "O canal onde esse comando for usado será utilizado para enviar avisos de players.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.SetServerUpdatesChannel(data.GS.ID, data.CS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaToggleDeaths = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaToggleDeaths",
	Aliases:     []string{"ttd", "senddeaths", "sd"},
	Description: "Determina se o bot enviará notícias de mortes de players ou não",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.ToggleDeaths(data.GS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaToggleUpdates = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaToggleUpdates",
	Aliases:     []string{"ttu", "sendupdates", "su"},
	Description: "Determina se o bot enviará notícias de players ou não",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Msg.Author.ID != data.GS.Guild.OwnerID {
			return "Apenas o dono do servidor pode usar esse comando.", nil
		}

		out, err := tibia.ToggleUpdates(data.GS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaGetWorld = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaGetWorld",
	Aliases:     []string{"tgw"},
	Description: "Retorna o mundo deste servidor.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.GetServerWorld(data.GS.ID, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaGetGuild = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "TibiaGetGuild",
	Aliases:     []string{"tgg"},
	Description: "Retorna a guild deste servidor.",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.GetServerGuild(data.GS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaCreateNewsFeed = &commands.YAGCommand{
	CmdCategory:         commands.CategoryTibia,
	Name:                "CreateNewsFeed",
	Aliases:             []string{"cnf"},
	Description:         "Cria o news feed.",
	RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
	Arguments: []*dcmd.ArgDef{
		{Name: "Canal para enviar as notícias", Type: dcmd.Channel},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		cID := data.Msg.ChannelID
		if data.Args[0].Value != nil {
			cID = data.Args[0].Value.(*dstate.ChannelState).ID
		}

		out, err := tibia.CreateNewsFeed(data.GS.ID, cID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaEnableNewsFeed = &commands.YAGCommand{
	CmdCategory:         commands.CategoryTibia,
	Name:                "EnableNewsFeed",
	Aliases:             []string{"enf"},
	Description:         "Habilita o news feed.",
	RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.EnableNewsFeed(data.GS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaDisableNewsFeed = &commands.YAGCommand{
	CmdCategory:         commands.CategoryTibia,
	Name:                "DisableNewsFeed",
	Aliases:             []string{"dnf"},
	Description:         "Desabilita o news feed.",
	RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.DisableNewsFeed(data.GS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TibiaChangeNewsFeedChannel = &commands.YAGCommand{
	CmdCategory:         commands.CategoryTibia,
	Name:                "ChangeNewsFeedChannel",
	Aliases:             []string{"cnfc"},
	Description:         "Troca o canal que as noticias do Tibia são enviadas.",
	RequireDiscordPerms: []int64{discordgo.PermissionManageServer},
	Arguments: []*dcmd.ArgDef{
		{Name: "Canal para enviar as notícias", Type: dcmd.Channel},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		cID := data.Msg.ChannelID
		if data.Args[0].Value != nil {
			cID = data.Args[0].Value.(*dstate.ChannelState).ID
		}

		out, err := tibia.ChangeNewsFeedChannel(data.GS.ID, cID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}
