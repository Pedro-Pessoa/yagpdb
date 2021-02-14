package info

import (
	"fmt"
	"math/rand"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Info",
	Description: "Responds with bot information",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		URL := common.ConfHost.GetString()
		embedPT := &discordgo.MessageEmbed{
			Title:       "TID BOT",
			Description: fmt.Sprintf("Tid bot é um bot de discord multi-propósito.\nO código original do bot foi escrito pelo Jonas747#0001(105487308693757952), e depois o Pedro Pessoa #8177(157321213641752576) bifurcou e modificou bastante o projeto criando esse bot.\n\nEsse bot foca em ser configurável, por isso ele é bem avançado.\n\nEsse bot pode performar uma grande gama de funções multi-propósito(Youtube/Twitter/Reddit feeds, vários comandos, utilidades de moderação, moderador automático, custom commands, cargo automático, menus de cargo, etc) e ele pode ser configurado através do painel de controle.\n\nO motivo para bifurcar do YAG foi para adicionar mais funcionalidades ao bot.\n\nTodas as mudanças podem ser encontradas aqui: <https://%s/manager/%d/changes\n\nPainel de controle: <https://%s/manage>", URL, data.Msg.GuildID, URL),
		}

		embedEN := &discordgo.MessageEmbed{
			Title:       "TID BOT",
			Color:       int(rand.Int63n(16777215)),
			Description: fmt.Sprintf("Tid bot is a multipuporse discord bot.\nThe original source code for the bot was written by Jonas747#0001(ID: 105487308693757952), and then Pedro Pessoa #8177(ID: 157321213641752576) forked and heavily modified it creating this bot.\n\nIt focus on being configurable, therefore it is rather advanced.\n\nIt can perform a range of general purpose functionality (Youtube/Twitter/Reddit feeds, various commands, moderation utilities, automoderator functionality, custom commands, autorole, role menus, and so on) and it's configured through a web control panel.\n\nThe reason for forkig out of YAG's was to add more functionalities to the bot.\n\nAll the changes can be found here: <https://%s/manager/%d/changes\n\nControl panel: <https://%s/manage>", URL, data.Msg.GuildID, URL),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Version",
					Value:  common.CurrentVersion,
					Inline: true,
				},
			},
		}

		serverLang := "PT"

		if serverLang == "PT" {
			return embedPT, nil
		}

		return embedEN, nil
	},
}
