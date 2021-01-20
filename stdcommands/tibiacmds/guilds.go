package tibiacmds

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/stdcommands/util"
	"github.com/jonas747/yagpdb/tibia"
)

var SpecificGuildCommand = &commands.YAGCommand{
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
		descOut := ""
		if len(desc) > 1700 {
			split := strings.Split(desc, " ")
			for i := range split {
				if len(descOut) < 1700 {
					descOut += fmt.Sprintf("%s, ", split[i])
				} else {
					descOut += "..."
					break
				}
			}
		} else {
			descOut = desc
		}

		embed := &discordgo.MessageEmbed{
			Title:       guild.Name,
			Color:       int(rand.Int63n(16777215)),
			Description: descOut,
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
