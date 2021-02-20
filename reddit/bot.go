package reddit

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/reddit/models"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var _ bot.RemoveGuildHandler = (*Plugin)(nil)

func (p *Plugin) RemoveGuild(g int64) error {
	_, err := models.RedditFeeds(models.RedditFeedWhere.GuildID.EQ(g)).UpdateAllG(context.Background(), models.M{
		"disabled": true,
	})
	if err != nil {
		return errors.WrapIf(err, "failed removing reddit feeds")
	}

	return nil
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, &commands.TIDCommand{
		CmdCategory:          commands.CategoryDebug,
		HideFromCommandsPage: true,
		Name:                 "testreddit",
		Description:          "Tests the reddit feeds in this server by checking the specified post",
		HideFromHelp:         true,
		RequiredArgs:         1,
		IsModCmd:             true,
		Arguments: []*dcmd.ArgDef{
			{Name: "post-id", Type: dcmd.String},
		},
		RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
			pID := data.Args[0].Str()
			if !strings.HasPrefix(pID, "t3_") {
				pID = "t3_" + pID
			}

			resp, err := p.redditClient.LinksInfo([]string{pID})
			if err != nil {
				return nil, err
			}

			if len(resp) < 1 {
				return "Unknown post", nil
			}

			handlerSlow := &PostHandlerImpl{
				Slow:        true,
				ratelimiter: NewRatelimiter(),
			}

			handlerFast := &PostHandlerImpl{
				Slow:        false,
				ratelimiter: NewRatelimiter(),
			}

			err1 := handlerSlow.handlePost(resp[0], data.GS.ID)
			err2 := handlerFast.handlePost(resp[0], data.GS.ID)

			return fmt.Sprintf("SlowErr: `%v`, fastErr: `%v`", err1, err2), nil
		}),
	})
}
