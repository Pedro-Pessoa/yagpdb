package topgames

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

var Command = &commands.TIDCommand{
	Cooldown:     5,
	CmdCategory:  commands.CategoryDebug,
	Name:         "topgames",
	Description:  "Shows the top games on this server",
	HideFromHelp: true,
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "all"},
	},
	RunFunc: cmdFuncTopCommands,
}

func cmdFuncTopCommands(data *dcmd.Data) (interface{}, error) {
	allSwitch := data.Switch("all")
	var all bool
	if allSwitch != nil && allSwitch.Bool() {
		all = true

		if admin, err := bot.IsBotAdmin(data.Msg.Author.ID); !admin || err != nil {
			if err != nil {
				return nil, err
			}

			return "Only bot admins can check top games of all servers", nil
		}
	}

	// do it in 2 passes for speedy accumulation of data
	fastResult := make(map[string]int)

	if all {
		guilds := bot.State.GuildsSlice(true)
		for _, g := range guilds {
			checkGuild(fastResult, g)
		}
	} else {
		checkGuild(fastResult, data.GS)
	}

	// then we convert and sort it
	fullResult := make([]*TopGameResult, 0, len(fastResult))
	for k, v := range fastResult {
		fullResult = append(fullResult, &TopGameResult{
			Game:  k,
			Count: v,
		})
	}

	sort.Slice(fullResult, func(i, j int) bool {
		return fullResult[i].Count > fullResult[j].Count
	})

	// display it
	var out strings.Builder
	out.WriteString("```\nTop games being played currently\n#    Count -  Game\n")
	for k, result := range fullResult {
		out.WriteString(fmt.Sprintf("#%02s: %5s - %s\n", strconv.Itoa(k+1), strconv.Itoa(result.Count), result.Game))
		if k >= 20 {
			break
		}
	}
	out.WriteString("\n```")

	return out.String(), nil
}

func checkGuild(dst map[string]int, gs *dstate.GuildState) {
	gs.RLock()
	defer gs.RUnlock()

	for _, ms := range gs.Members {
		if !ms.PresenceSet || ms.PresenceActivities == nil || len(ms.PresenceActivities) == 0 {
			continue
		}

		if ms.Bot {
			continue
		}

		for _, p := range ms.PresenceActivities {
			if p != nil {
				name := p.Name
				dst[name]++
			}
		}
	}
}

type TopGameResult struct {
	Game  string
	Count int
}
