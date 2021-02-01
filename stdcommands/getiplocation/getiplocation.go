// From PAGST

package getiplocation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

type ipAPIJSON struct {
	Status        string
	Message       string
	Continent     string
	ContinentCode string
	Country       string
	CountryCode   string
	Region        string
	RegionName    string
	City          string
	District      string
	Zip           string
	Lat           float64
	Lon           float64
	Timezone      string
	Offset        int64
	Currency      string
	Isp           string
	Org           string
	As            string
	Asname        string
	Reverse       string
	Mobile        bool
	Proxy         bool
	Hosting       bool
	Query         string
}

var Command = &commands.TIDCommand{
	CmdCategory:  commands.CategoryTool,
	Name:         "GetIPLocation",
	Aliases:      []string{"geoloc", "getiploc", "iploc"},
	Description:  "Queries IP Geolocation API on given IP-address or domain",
	RunInDM:      true,
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "IP-address or domain", Type: dcmd.String},
	},

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		var ipAPIHost = "ip-api.com" //has 45 requests per minute for free account and no HTTPS
		var queryType = "json"
		var ipJSON ipAPIJSON

		//we make the queryURL here
		ipArg := data.Args[0].Str()
		queryURL := fmt.Sprintf("http://%s/%s/%s", ipAPIHost, queryType, ipArg)

		//let's get that API data
		body, err := getData(queryURL)
		if err != nil {
			return nil, err
		}

		queryErr := json.Unmarshal([]byte(body), &ipJSON)
		if queryErr != nil {
			return nil, queryErr
		}

		if ipJSON.Status == "fail" {
			return nil, commands.NewPublicError("Cannot fetch IP-location from given data: ", ipArg)
		}

		if ipJSON.Org == "" {
			ipJSON.Org = "-"
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Query: %s", ipJSON.Query),
			URL:         fmt.Sprintf("https://%s/%s", ipAPIHost, ipJSON.Query),
			Description: fmt.Sprintf("**Country:**\n%s (%s)\n\n**City/Region:**\n%s, %s\n\n**ISP/ORG:**\n%s; %s", ipJSON.Country, ipJSON.CountryCode, ipJSON.City, ipJSON.RegionName, ipJSON.Isp, ipJSON.Org),
			Color:       int(rand.Int63n(16777215)),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://ip-api.com/docs/static/logo.png",
			},
		}

		return embed, nil
	},
}

func getData(query string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(query)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, commands.NewPublicError("Cannot fetch IP-location. Try again later.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
