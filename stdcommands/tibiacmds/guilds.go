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

var SpecificGuildCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "Guild",
	Description: "Retorna informações da guild especificada.",
	Cooldown:    5,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome da Guild", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		guild, err := tibia.GetTibiaSpecificGuild(data.Args[0].Str())
		if err != nil {
			return fmt.Sprintln(err), err
		}

		desc := guild.Description
		var descOut strings.Builder
		if len(desc) > 1700 {
			split := strings.Split(desc, " ")
			for _, s := range split {
				if len(descOut.String()) < 1700 {
					descOut.WriteString(s + " ")
				} else {
					descOut.WriteString("...")
					break
				}
			}
		} else {
			descOut.WriteString(desc)
		}

		embed := &discordgo.MessageEmbed{
			Title:       guild.Name,
			Color:       int(rand.Int63n(16777215)),
			Description: descOut.String(),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Número de membros",
					Value:  strconv.Itoa(guild.MemberCount),
					Inline: true,
				},
				{
					Name:   "Mundo",
					Value:  guild.World,
					Inline: true,
				},
				{
					Name:   "Guild Hall",
					Value:  guild.GuildHall,
					Inline: true,
				},
				{
					Name:   "Está em Guerra?",
					Value:  guild.War,
					Inline: true,
				},
			},
		}

		return embed, nil
	},
}
