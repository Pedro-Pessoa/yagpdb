package tibiacmds

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/tibia"
)

var MainCharCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "Char",
	Description: "Retorna informações do personagem especificado.",
	Cooldown:    5,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Source == dcmd.DMSource {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		char, err := tibia.GetTibiaChar(data.Args[0].Str(), true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		comentario := char.Comment
		comentarioOut := ""
		if len(comentario) > 1700 {
			split := strings.Split(comentario, " ")
			for i := range split {
				if len(comentarioOut) < 1700 {
					comentarioOut += fmt.Sprintf("%s, ", split[i])
				} else {
					comentarioOut += "..."
					break
				}
			}
		} else {
			comentarioOut = comentario
		}
		re := regexp.MustCompile(` `)
		linkname := re.ReplaceAllString(char.Name, `\+`)
		link := fmt.Sprintf("https://www.tibia.com/community/?subtopic=characters&name=%s", linkname)
		comentarioOut = fmt.Sprintf("%s\n\n[Perfil do char](%s)", comentarioOut, link)

		embed := &discordgo.MessageEmbed{
			Title:       char.Name,
			Color:       int(rand.Int63n(16777215)),
			Description: comentarioOut,
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

var DeathsCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "Mortes",
	Description: "Retorna as mortes recentes do personagem especificado.",
	Cooldown:    5,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Source == dcmd.DMSource {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		char, err := tibia.GetTibiaChar(data.Args[0].Str(), false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		mortes := char.Deaths
		deaths := "\n"
		motivo := ""

		if len(mortes) >= 1 {
			re := regexp.MustCompile(`Died by a`)
			for _, v := range mortes {
				if len(deaths) < 1800 {
					checkKillByMonster := re.MatchString(v.Reason)
					if checkKillByMonster {
						deaths += fmt.Sprintf("**Data**: %s\n**Level**: %d\n**Motivo**: %s\n\n", v.Date, v.Level, v.Reason)
					} else {
						split := strings.Split(v.Reason, ",")
						for i := range split {
							if len(motivo) < 150 {
								motivo += fmt.Sprintf("%s, ", split[i])
							} else {
								motivo += "e outros."
								break
							}
						}
						re := regexp.MustCompile(`,\s*\z`)
						motivo = re.ReplaceAllString(motivo, ".")
						deaths += fmt.Sprintf("**Data**: %s\n**Level**: %d\n**Motivo**: %s\n\n", v.Date, v.Level, motivo)
						motivo = ""
					}
				} else {
					deaths += "... entre outras ..."
					break
				}
			}
		} else {
			deaths = "Sem mortes recentes."
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Mortes recentes de %s", char.Name),
			Description: deaths,
			Color:       int(rand.Int63n(16777215)),
		}

		return embed, nil

	},
}

var CheckOnlineCommand = &commands.YAGCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "CheckOnline",
	Description: "Mostra quem está online no mundo especificado.",
	Aliases:     []string{"co"},
	Cooldown:    10,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Mundo", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Source == dcmd.DMSource {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		mundo, name, err := tibia.CheckOnline(data.Args[0].Str())
		if err != nil {
			return fmt.Sprintln(err), err
		}

		desc := ""
		if len(mundo) > 0 {
			for _, v := range mundo {
				if len(desc) < 1700 {
					desc += fmt.Sprintf("%s, ", v.Name)
				} else {
					desc += "e outros."
					break
				}
			}
			re := regexp.MustCompile(`,\s*\z`)
			url := fmt.Sprintf("https://www.tibia.com/community/?subtopic=worlds&world=%s", *name)
			desc = fmt.Sprintf("%s\n\n[Veja todas os players online](%s)", re.ReplaceAllString(desc, "."), url)
		} else {
			desc = "Nenhum player online."
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Players online em %s", *name),
			Description: desc,
			Color:       int(rand.Int63n(16777215)),
		}

		return embed, nil

	},
}
