package tibiacmds

import (
	"fmt"
	"math/rand"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
	"github.com/Pedro-Pessoa/tidbot/tibia"
)

var NewsCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "News",
	Aliases:     []string{"noticia"},
	Description: "Última noticia do tibia, ou alguma específica.",
	Cooldown:    10,
	Arguments: []*dcmd.ArgDef{
		{Name: "ID da notícia", Type: dcmd.Int},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		var news *tibia.InternalNews
		var err error

		if data.Args[0].Value != nil {
			news, err = tibia.GetTibiaNews(data.Args[0].Int())
			if err != nil {
				return fmt.Sprintln(err), err
			}
		} else {
			news, err = tibia.GetTibiaNews()
			if err != nil {
				return fmt.Sprintln(err), err
			}
		}

		embed := &discordgo.MessageEmbed{
			Title:       news.Title,
			Color:       int(rand.Int63n(16777215)),
			Description: fmt.Sprintf("%s\n[Clique para ver mais](%s)", news.ShortDescription, news.URL),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("ID: %d\nData: %s", news.ID, news.Date),
			},
		}

		return embed, nil

	},
}

var NewsTickerCommand = &commands.TIDCommand{
	CmdCategory: commands.CategoryTibia,
	Name:        "NewsTicker",
	Description: "Último newsticker do tibia.",
	Cooldown:    10,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		news, err := tibia.GetTibiaNewsticker()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		embed := &discordgo.MessageEmbed{
			Title:       news.Title,
			Color:       int(rand.Int63n(16777215)),
			Description: fmt.Sprintf("%s\n[Clique para ver mais](%s)", news.ShortDescription, news.URL),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Notícia mais recente do Tibia. | ID: %d\nData: %s", news.ID, news.Date),
			},
		}

		return embed, nil

	},
}
