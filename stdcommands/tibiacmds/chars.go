package tibiacmds

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
	"github.com/Pedro-Pessoa/tidbot/tibia"
)

var MainCharCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "Char",
	Description: "Retorna informações do personagem especificado.",
	Cooldown:    5,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		char, err := tibia.GetTibiaChar(data.Args[0].Str(), true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		comentario := char.Comment
		var comentarioOut strings.Builder
		if len(comentario) > 1700 {
			split := strings.Split(comentario, " ")
			for i := range split {
				if len(comentarioOut.String()) < 1700 {
					comentarioOut.WriteString(split[i] + ", ")
				} else {
					comentarioOut.WriteString("...")
					break
				}
			}
		} else {
			comentarioOut.WriteString(comentario)
		}

		linkname := strings.ReplaceAll(char.Name, " ", "+")
		link := "https://www.tibia.com/community/?subtopic=characters&name=" + linkname
		comentarioOut.WriteString("\n\n[Perfil do char](" + link + ")")

		embed := &discordgo.MessageEmbed{
			Title:       char.Name,
			Color:       int(rand.Int63n(16777215)),
			Description: comentarioOut.String(),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Level",
					Value:  strconv.Itoa(char.Level),
					Inline: true,
				},
				{
					Name:   "Mundo",
					Value:  char.World,
					Inline: true,
				},
				{
					Name:   "Vocação",
					Value:  char.Vocation,
					Inline: true,
				},
				{
					Name:   "Templo",
					Value:  char.Residence,
					Inline: true,
				},
				{
					Name:   "Status",
					Value:  char.AccountStatus,
					Inline: true,
				},
				{
					Name:   "On/Off",
					Value:  strings.Title(char.Status),
					Inline: true,
				},
				{
					Name:   "Lealdade",
					Value:  char.Loyalty,
					Inline: true,
				},
				{
					Name:   "Pontos de Achievement",
					Value:  strconv.Itoa(char.AchievementPoints),
					Inline: true,
				},
				{
					Name:   "Gênero",
					Value:  strings.Title(char.Sex),
					Inline: true,
				},
				{
					Name:   "Guild",
					Value:  char.Guild,
					Inline: true,
				},
			},
		}

		if char.Rank != "Sem guild" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Cargo na Guild",
				Value:  char.Rank,
				Inline: true,
			})
		}

		if char.House != "Nenhuma" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Casa",
				Value:  char.House,
				Inline: true,
			})
		}

		if char.Married != "Ninguém" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Casado",
				Value:  char.Married,
				Inline: true,
			})
		}

		if char.CreatedAt != "Data escondida" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Criado",
				Value:  char.CreatedAt,
				Inline: true,
			})
		}

		return embed, nil
	},
}

var DeathsCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "Mortes",
	Description: "Retorna as mortes recentes do personagem especificado.",
	Cooldown:    5,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		char, err := tibia.GetTibiaChar(data.Args[0].Str(), false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		mortes := char.Deaths
		var deaths strings.Builder

		for _, v := range mortes {
			if len(deaths.String()) < 1800 {
				if strings.Contains(v.Reason, "Died by a") { // check if the char was killed by a monster
					deaths.WriteString("**Data**: " + v.Date + "\n**Level**: " + strconv.Itoa(v.Level) + "\n**Motivo**: " + v.Reason + "\n\n")
					continue
				}

				split := strings.Split(v.Reason, ",")
				var motivo strings.Builder
				for i, s := range split {
					if i == 0 {
						motivo.WriteString(s)
						continue
					}

					if len(motivo.String()) < 150 {
						motivo.WriteString(", " + s)
					} else {
						motivo.WriteString(" e outros...")
						break
					}
				}

				deaths.WriteString("**Data**: " + v.Date + "\n**Level**: " + strconv.Itoa(v.Level) + "\n**Motivo**: " + motivo.String() + "\n\n")
			} else {
				deaths.WriteString("... entre outras ...")
				break
			}
		}

		if len(mortes) == 0 {
			deaths.WriteString("Sem mortes recentes.")
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Mortes recentes de " + char.Name,
			Description: deaths.String(),
			Color:       int(rand.Int63n(16777215)),
		}

		return embed, nil
	},
}

var CheckOnlineCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "CheckOnline",
	Description: "Mostra quem está online no mundo especificado.",
	Aliases:     []string{"co"},
	Cooldown:    10,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Mundo", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		mundo, name, err := tibia.CheckOnline(data.Args[0].Str())
		if err != nil {
			return fmt.Sprintln(err), err
		}

		var desc strings.Builder
		if len(mundo) > 0 {
			for i, v := range mundo {
				if i == 0 {
					desc.WriteString(v.Name)
					continue
				}

				if len(desc.String()) < 1700 {
					desc.WriteString(", " + v.Name)
				} else {
					desc.WriteString(" e outros...")
					break
				}
			}

			url := "https://www.tibia.com/community/?subtopic=worlds&world=" + name
			desc.WriteString("\n\n[Veja todos os players online](" + url + ")")
		} else {
			desc.WriteString("Nenhum player online.")
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Players online em " + name,
			Description: desc.String(),
			Color:       int(rand.Int63n(16777215)),
		}

		return embed, nil
	},
}
