package topevents

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	Cooldown:     2,
	CmdCategory:  commands.CategoryDebug,
	Name:         "topevents",
	Description:  "Shows gateway event processing stats for all or one shard",
	HideFromHelp: true,
	Arguments: []*dcmd.ArgDef{
		{Name: "shard", Type: dcmd.Int},
	},
	RunFunc: cmdFuncTopEvents,
}

func cmdFuncTopEvents(data *dcmd.Data) (interface{}, error) {
	shardsTotal, lastPeriod := bot.EventLogger.GetStats()

	sortable := make([]*DiscordEvtEntry, len(eventsystem.AllDiscordEvents))
	for i := range sortable {
		sortable[i] = &DiscordEvtEntry{
			Name: eventsystem.AllDiscordEvents[i].String(),
		}
	}

	for i := range shardsTotal {
		if data.Args[0].Value != nil && data.Args[0].Int() != i {
			continue
		}

		for de, j := range eventsystem.AllDiscordEvents {
			sortable[de].Total += shardsTotal[i][j]
			sortable[de].PerSecond += float64(lastPeriod[i][j]) / bot.EventLoggerPeriodDuration.Seconds()
		}
	}

	sort.Sort(DiscordEvtEntrySortable(sortable))

	var out strings.Builder
	out.WriteString("Total event stats across all shards:\n")

	if data.Args[0].Value != nil {
		out.WriteString("Stats for shard " + strconv.Itoa(data.Args[0].Int()) + ":\n")
	}

	out.WriteString("```\n#     Total  -   /s  - Event\n")

	var sum int64
	var sumPerSecond float64

	for k, entry := range sortable {
		out.WriteString(fmt.Sprintf("#%-2s: %7s - %5.1f - %s\n", strconv.Itoa(k+1), strconv.FormatInt(entry.Total, 10), entry.PerSecond, entry.Name))
		sum += entry.Total
		sumPerSecond += entry.PerSecond
	}

	out.WriteString(fmt.Sprintf("\nTotal: %s, Events per second: %.1f", strconv.FormatInt(sum, 10), sumPerSecond))
	out.WriteString("\n```")

	return out.String(), nil
}

type DiscordEvtEntry struct {
	Name      string
	Total     int64
	PerSecond float64
}

type DiscordEvtEntrySortable []*DiscordEvtEntry

func (d DiscordEvtEntrySortable) Len() int {
	return len(d)
}

func (d DiscordEvtEntrySortable) Less(i, j int) bool {
	return d[i].Total > d[j].Total
}

func (d DiscordEvtEntrySortable) Swap(i, j int) {
	temp := d[i]
	d[i] = d[j]
	d[j] = temp
}
