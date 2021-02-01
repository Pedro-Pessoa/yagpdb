package xkcd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

type Xkcd struct {
	Month      string
	Num        int64
	Link       string
	Year       string
	News       string
	SafeTitle  string
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

var XkcdHost = "https://xkcd.com/"
var XkcdJson = "info.0.json"

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Xkcd",
	Description: "An xkcd comic, by default returns random comic strip",
	Arguments: []*dcmd.ArgDef{
		{Name: "Comic number", Type: dcmd.Int},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "l", Name: "Latest comic"},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		//first query to get latest number
		latest := false
		xkcd, err := getComic()
		if err != nil {
			return "Something happened whilst getting the comic!", err
		}

		xkcdNum := rand.Int63n(xkcd.Num) + 1

		//latest comic strip flag, already got that data
		if data.Switches["l"].Value != nil && data.Switches["l"].Value.(bool) {
			latest = true
		}

		//specific comic strip number
		if data.Args[0].Value != nil {
			n := data.Args[0].Int64()
			if n >= 1 && n <= xkcd.Num {
				xkcdNum = n
			} else {
				return fmt.Sprintf("There's no comic numbered %d, current range is 1-%d", n, xkcd.Num), nil
			}
		}

		//if no latest flag is set, fetches a comic by number
		if !latest {
			xkcd, err = getComic(xkcdNum)
			if err != nil {
				return "Something happened whilst getting the comic!", err
			}
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("#%d: %s", xkcd.Num, xkcd.Title),
			Description: fmt.Sprintf("[%s](%s%d/)", xkcd.Alt, XkcdHost, xkcd.Num),
			Color:       int(rand.Int63n(16777215)),
			Image: &discordgo.MessageEmbedImage{
				URL: xkcd.Img,
			},
		}

		return embed, nil
	},
}

func getComic(number ...int64) (*Xkcd, error) {
	xkcd := Xkcd{}
	queryUrl := XkcdHost + XkcdJson

	if len(number) >= 1 {
		queryUrl = fmt.Sprintf(XkcdHost+"%d/"+XkcdJson, number[0])
	}

	resp, err := http.DefaultClient.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&xkcd)
	if err != nil {
		return nil, err
	}

	return &xkcd, nil
}
