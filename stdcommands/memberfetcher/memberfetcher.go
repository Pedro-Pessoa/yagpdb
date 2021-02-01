package memberfetcher

import (
	"fmt"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryDebug,
	Name:        "MemberFetcher",
	Aliases:     []string{"memfetch"},
	Description: "Shows the current status of the member fetcher",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		fetching, notFetching := bot.MemberFetcher.Status()
		return fmt.Sprintf("Fetching: `%d`, Not fetching: `%d`", fetching, notFetching), nil
	},
}
