package ping

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory:     commands.CategoryDebug,
	Name:            "Ping",
	Description:     "Shows the latency from the bot to the discord servers.",
	LongDescription: "Note that high latencies can be the fault of ratelimits and the bot itself, it's not a absolute metric.",

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return fmt.Sprintf(":PONG;%d", time.Now().UnixNano()), nil
	},
}

func HandleMessageCreate(evt *eventsystem.EventData) {
	m := evt.MessageCreate()

	bUser := common.BotUser
	if bUser == nil {
		return
	}

	if bUser.ID != m.Author.ID {
		return
	}

	// ping pong
	split := strings.Split(m.Content, ";")
	if split[0] != ":PONG" || len(split) < 2 {
		return
	}

	parsed, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return
	}

	taken := time.Duration(time.Now().UnixNano() - parsed)

	started := time.Now()
	_, _ = common.BotSession.ChannelMessageEdit(m.ChannelID, m.ID, "Gateway (http send -> gateway receive time): "+taken.String())
	httpPing := time.Since(started)

	_, _ = common.BotSession.ChannelMessageEdit(m.ChannelID, m.ID, "HTTP API (Edit Msg): "+httpPing.String()+"\nGateway: "+taken.String())
}
