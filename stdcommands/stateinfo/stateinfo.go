package stateinfo

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
)

var Command = &commands.YAGCommand{
	Cooldown:     2,
	CmdCategory:  commands.CategoryDebug,
	Name:         "stateinfo",
	Description:  "Responds with state debug info",
	HideFromHelp: true,
	RunFunc:      cmdFuncStateInfo,
}

func cmdFuncStateInfo(data *dcmd.Data) (interface{}, error) {
	totalGuilds := 0
	totalMembers := 0
	guildChannel := 0
	totalMessages := 0

	state := bot.State
	state.RLock()
	totalChannels := len(state.Channels)
	totalGuilds = len(state.Guilds)
	gCop := state.GuildsSlice(false)
	state.RUnlock()

	for _, g := range gCop {
		g.RLock()

		guildChannel += len(g.Channels)
		totalMembers += len(g.Members)

		for _, cState := range g.Channels {
			totalMessages += len(cState.Messages)
		}
		g.RUnlock()
	}

	stats := bot.State.StateStats()

	embed := &discordgo.MessageEmbed{
		Title: "State size",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Guilds", Value: fmt.Sprint(totalGuilds), Inline: true},
			{Name: "Members", Value: fmt.Sprintf("%d", totalMembers), Inline: true},
			{Name: "Messages", Value: fmt.Sprintf("%d", totalMessages), Inline: true},
			{Name: "Guild Channels", Value: fmt.Sprintf("%d", guildChannel), Inline: true},
			{Name: "Total Channels", Value: fmt.Sprintf("%d", totalChannels), Inline: true},
			{Name: "Cache Hits/Misses", Value: fmt.Sprintf("%d - %d", stats.CacheHits, stats.CacheMisses), Inline: true},
			{Name: "Members evicted total", Value: fmt.Sprintf("%d", stats.MembersRemovedTotal), Inline: true},
			{Name: "Cache evicted total", Value: fmt.Sprintf("%d", stats.UserCachceEvictedTotal), Inline: true},
			{Name: "Messages removed total", Value: fmt.Sprintf("%d", stats.MessagesRemovedTotal), Inline: true},
		},
	}

	return embed, nil
}
