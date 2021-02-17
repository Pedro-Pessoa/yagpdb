package tibiacmds

import (
	"fmt"
	"math/rand"
	"strconv"

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

		return embedBuilder(news), nil

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

		return embedBuilder(news), nil
	},
}

func embedBuilder(news *tibia.InternalNews) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       news.Title,
		Color:       int(rand.Int63n(16777215)),
		Description: news.ShortDescription + "\n[Clique aqui para ver mais](" + news.URL + ")",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ID: " + strconv.Itoa(news.ID) + "\nData: " + news.Date,
		},
	}
}
