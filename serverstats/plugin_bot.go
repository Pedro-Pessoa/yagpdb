package serverstats

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/pubsub"
	"github.com/jonas747/yagpdb/serverstats/messagestatscollector"
	"github.com/jonas747/yagpdb/web"
	"github.com/mediocregopher/radix/v3"
)

func MarkGuildAsToBeChecked(guildID int64) {
	common.RedisPool.Do(radix.FlatCmd(nil, "SADD", "serverstats_active_guilds", guildID))
}

var (
	_                   bot.BotInitHandler       = (*Plugin)(nil)
	_                   commands.CommandProvider = (*Plugin)(nil)
	msgStatsCollector   *messagestatscollector.Collector
	memberSatatsUpdater *serverMemberStatsUpdater
)

func (p *Plugin) BotInit() {
	msgStatsCollector = messagestatscollector.NewCollector(logger, time.Minute*5)
	memberSatatsUpdater = newServerMemberStatsUpdater()
	go memberSatatsUpdater.run()

	pubsub.AddHandler("server_stats_invalidate_cache", func(evt *pubsub.Event) {
		gs := bot.State.Guild(true, evt.TargetGuildInt)
		if gs != nil {
			gs.UserCacheDel(CacheKeyConfig)
		}
	}, nil)

	if !confDeprecated.GetBool() {
		eventsystem.AddHandlerAsyncLastLegacy(p, handleUpdateMemberStats, eventsystem.EventGuildMemberAdd, eventsystem.EventGuildMemberRemove, eventsystem.EventGuildCreate)
		eventsystem.AddHandlerAsyncLast(p, eventsystem.RequireCSMW(HandleMessageCreate), eventsystem.EventMessageCreate)
		go p.runOnlineUpdater()
	} else {
		logger.Info("Not enabling server stats collecting due to deprecation flag being set")
	}
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, &commands.YAGCommand{
		CustomEnabled: true,
		CmdCategory:   commands.CategoryTool,
		Cooldown:      5,
		Name:          "Stats",
		Description:   "Shows server stats (if public stats are enabled)",
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			config, err := GetConfig(data.Context(), data.GS.ID)
			if err != nil {
				return nil, errors.WithMessage(err, "getconfig")
			}

			if !config.Public {
				return fmt.Sprintf("Stats are set to private on this server, this can be changed in the control panel on <https://%s>", common.ConfHost.GetString()), nil
			}

			stats, err := RetrieveDailyStats(time.Now(), data.GS.ID)
			if err != nil {
				return nil, errors.WithMessage(err, "retrievefullstats")
			}

			total := int64(0)
			for _, c := range stats.ChannelMessages {
				total += c.Count
			}

			embed := &discordgo.MessageEmbed{
				Title:       "Server stats",
				Description: fmt.Sprintf("[Click here to open in browser](%s/public/%d/stats)", web.BaseURL(), data.GS.ID),
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Members joined 24h", Value: fmt.Sprint(stats.JoinedDay), Inline: true},
					{Name: "Members Left 24h", Value: fmt.Sprint(stats.LeftDay), Inline: true},
					{Name: "Total Messages 24h", Value: fmt.Sprint(total), Inline: true},
					{Name: "Members Online", Value: fmt.Sprint(stats.Online), Inline: true},
					{Name: "Total Members", Value: fmt.Sprint(stats.TotalMembers), Inline: true},
				},
			}

			return embed, nil
		},
	})
}

func handleUpdateMemberStats(evt *eventsystem.EventData) {
	select {
	case memberSatatsUpdater.incoming <- evt:
	default:
		go func() {
			memberSatatsUpdater.incoming <- evt
		}()
	}
}

func HandleMessageCreate(evt *eventsystem.EventData) (retry bool, err error) {

	m := evt.MessageCreate()
	if m.GuildID == 0 || m.Author == nil || m.Author.Bot {
		return // private channel
	}

	channel := evt.CS()

	config, err := BotCachedFetchGuildConfig(evt.Context(), channel.Guild)
	if err != nil {
		return true, errors.WithStackIf(err)
	}

	if common.ContainsInt64Slice(config.ParsedChannels, channel.ID) {
		return false, nil
	}

	msgStatsCollector.MsgEvtChan <- m.Message
	return false, nil
}

type CacheKey int

const (
	CacheKeyConfig CacheKey = iota
)

func BotCachedFetchGuildConfig(ctx context.Context, gs *dstate.GuildState) (*ServerStatsConfig, error) {
	v, err := gs.UserCacheFetch(CacheKeyConfig, func() (interface{}, error) {
		return GetConfig(ctx, gs.ID)
	})

	if err != nil {
		return nil, err
	}

	return v.(*ServerStatsConfig), nil
}

func keyOnlineMembers(year, day int) string {
	return "serverstats_online_members:" + strconv.Itoa(year) + ":" + strconv.Itoa(day)
}

func (p *Plugin) runOnlineUpdater() {
	time.Sleep(time.Minute * 1) // relieve startup preasure

	ticker := time.NewTicker(time.Second * 10)
	state := bot.State

	var guildsToCheck []*dstate.GuildState
	var i int
	var numToCheckPerRun int

	for {
		select {
		case <-ticker.C:
		}

		if len(guildsToCheck) < 0 || i >= len(guildsToCheck) {
			// Copy the list of guilds so that we dont need to keep the entire state locked

			i = 0
			guildsToCheck = state.GuildsSlice(true)

			// Hit each guild once per hour more or less
			numToCheckPerRun = len(guildsToCheck) / 250
			if numToCheckPerRun < 1 {
				numToCheckPerRun = 1
			}
		}

		started := time.Now()

		totalCounts := make(map[int64][2]int)

		checkedThisRound := 0
		for ; i < len(guildsToCheck) && checkedThisRound < numToCheckPerRun; i++ {
			g := guildsToCheck[i]
			online, total := p.checkGuildOnlineCount(g)

			totalCounts[g.ID] = [2]int{online, total}
			checkedThisRound++
		}

		t := time.Now()
		day := t.YearDay()
		year := t.Year()

		updateActions := make([]radix.CmdAction, 0, len(totalCounts)*2)

		for g, counts := range totalCounts {
			updateActions = append(updateActions, radix.FlatCmd(nil, "ZADD", keyTotalMembers(year, day), counts[1], g), radix.FlatCmd(nil, "ZADD", keyOnlineMembers(year, day), counts[0], g))
		}

		err := common.RedisPool.Do(radix.Pipeline(updateActions...))
		if err != nil {
			logger.WithError(err).Error("failed updating members period runner")
		}

		if time.Since(started) > time.Second {
			logger.Warnf("Tok %s to update online counts of %d guilds", time.Since(started).String(), checkedThisRound)
		}
	}
}

func (p *Plugin) checkGuildOnlineCount(guild *dstate.GuildState) (online int, total int) {

	guild.RLock()
	total = guild.Guild.MemberCount
	for _, v := range guild.Members {
		if v.PresenceSet && v.PresenceStatus != dstate.StatusOffline {
			online++
		}
	}
	guild.RUnlock()

	return online, total
}
