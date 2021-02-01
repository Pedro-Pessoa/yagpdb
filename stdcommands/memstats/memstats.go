package memstats

import (
	"bytes"
	"encoding/json"
	"runtime"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "memstats",
	Description:          ";))",
	HideFromHelp:         true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		buf, _ := json.Marshal(m)

		send := &discordgo.MessageSend{
			Content: "Memory stats",
			File: &discordgo.File{
				ContentType: "application/json",
				Name:        "memory_stats.json",
				Reader:      bytes.NewReader(buf),
			},
		}

		_, err := common.BotSession.ChannelMessageSendComplex(data.Msg.ChannelID, send)

		return nil, err
	}),
}
