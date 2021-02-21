package notifications

import (
	"fmt"
	"math/rand"
	"strings"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/analytics"
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/templates"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

var _ bot.BotInitHandler = (*Plugin)(nil)

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLast(p, HandleGuildMemberAdd, eventsystem.EventGuildMemberAdd)
	eventsystem.AddHandlerAsyncLast(p, HandleGuildMemberUpdate, eventsystem.EventGuildMemberUpdate)
	eventsystem.AddHandlerAsyncLast(p, HandleGuildMemberRemove, eventsystem.EventGuildMemberRemove)
	eventsystem.AddHandlerFirst(p, HandleChannelUpdate, eventsystem.EventChannelUpdate)
}

func HandleGuildMemberAdd(evtData *eventsystem.EventData) (retry bool, err error) {
	evt := evtData.GuildMemberAdd()

	config, err := GetConfig(evt.GuildID)
	if err != nil {
		return true, errors.WithStackIf(err)
	}

	if !config.JoinServerEnabled && !config.JoinDMEnabled {
		return
	}

	if (!config.JoinDMEnabled || evt.User.Bot) && !config.JoinServerEnabled {
		return
	}

	gs := bot.State.Guild(true, evt.GuildID)
	ms := dstate.MSFromDGoMember(gs, evt.Member)

	// Beware of the pyramid and its curses
	if config.JoinDMEnabled && !evt.User.Bot {
		if config.JoinDMWaitForPendingEnabled {
			return
		}

		cid, err := common.BotSession.UserChannelCreate(evt.User.ID)
		if err != nil {
			if bot.CheckDiscordErrRetry(err) {
				return true, errors.WithStackIf(err)
			}

			logger.WithError(err).WithField("user", evt.User.ID).Error("Failed retrieving user channel")
		} else {
			thinCState := &dstate.ChannelState{
				Owner: gs,
				Guild: gs,
				ID:    cid.ID,
				Name:  evt.User.Username,
				Type:  discordgo.ChannelTypeDM,
			}

			go analytics.RecordActiveUnit(gs.ID, &Plugin{}, "posted_join_server_msg")

			if sendTemplate(thinCState, config.JoinDMMsg, ms, "join dm", false) {
				return true, nil
			}
		}
	}

	if config.JoinServerEnabled && len(config.JoinServerMsgs) > 0 {
		if config.JoinServerWaitForPendingEnabled {
			return
		}

		channel := gs.Channel(true, config.JoinServerChannelInt())
		if channel == nil {
			return
		}

		go analytics.RecordActiveUnit(gs.ID, &Plugin{}, "posted_join_server_dm")

		chanMsg := config.JoinServerMsgs[rand.Intn(len(config.JoinServerMsgs))]
		if sendTemplate(channel, chanMsg, ms, "join server msg", config.CensorInvites) {
			return true, nil
		}
	}

	return false, nil
}

func HandleGuildMemberUpdate(evtData *eventsystem.EventData) (retry bool, err error) {
	if !bot.State.TrackBeforeStates {
		return
	}

	evt := evtData.GuildMemberUpdate()

	if evt.Pending {
		return
	}

	config, err := GetConfig(evt.GuildID)
	if err != nil {
		return true, errors.WithStackIf(err)
	}

	if !config.JoinDMWaitForPendingEnabled && !config.JoinServerWaitForPendingEnabled {
		return
	}

	if !config.JoinServerEnabled && !config.JoinDMEnabled {
		return
	}

	if (!config.JoinDMEnabled || evt.User.Bot) && !config.JoinServerEnabled {
		return
	}

	gs := bot.State.Guild(true, evt.GuildID)
	ms := dstate.MSFromDGoMember(gs, evt.Member)

	bot.State.BeforeStateLocker.Lock()
	beforeMS := bot.State.BeforeStateMemberMap[evt.Member.User.ID]
	bot.State.BeforeStateLocker.Unlock()

	if beforeMS == nil || !beforeMS.MemberState.Pending {
		// Msg probably already sent
		return
	}

	// Beware of the pyramid and its curses
	if config.JoinDMEnabled && !evt.User.Bot {
		cid, err := common.BotSession.UserChannelCreate(evt.User.ID)
		if err != nil {
			if bot.CheckDiscordErrRetry(err) {
				return true, errors.WithStackIf(err)
			}

			logger.WithError(err).WithField("user", evt.User.ID).Error("Failed retrieving user channel")
		} else {
			thinCState := &dstate.ChannelState{
				Owner: gs,
				Guild: gs,
				ID:    cid.ID,
				Name:  evt.User.Username,
				Type:  discordgo.ChannelTypeDM,
			}

			go analytics.RecordActiveUnit(gs.ID, &Plugin{}, "posted_join_server_msg")

			if sendTemplate(thinCState, config.JoinDMMsg, ms, "join dm", false) {
				return true, nil
			}
		}
	}

	if config.JoinServerEnabled && len(config.JoinServerMsgs) > 0 {
		channel := gs.Channel(true, config.JoinServerChannelInt())
		if channel == nil {
			return
		}

		go analytics.RecordActiveUnit(gs.ID, &Plugin{}, "posted_join_server_dm")

		chanMsg := config.JoinServerMsgs[rand.Intn(len(config.JoinServerMsgs))]
		if sendTemplate(channel, chanMsg, ms, "join server msg", config.CensorInvites) {
			return true, nil
		}
	}

	return false, nil
}

func HandleGuildMemberRemove(evt *eventsystem.EventData) (retry bool, err error) {
	memberRemove := evt.GuildMemberRemove()

	config, err := GetConfig(memberRemove.GuildID)
	if err != nil {
		return true, errors.WithStackIf(err)
	}

	if !config.LeaveEnabled || len(config.LeaveMsgs) == 0 {
		return
	}

	gs := bot.State.Guild(true, memberRemove.GuildID)
	if gs == nil {
		return
	}

	channel := gs.Channel(true, config.LeaveChannelInt())
	if channel == nil {
		return
	}

	ms := dstate.MSFromDGoMember(gs, memberRemove.Member)
	chanMsg := config.LeaveMsgs[rand.Intn(len(config.LeaveMsgs))]

	go analytics.RecordActiveUnit(gs.ID, &Plugin{}, "posted_leave_server_msg")

	if sendTemplate(channel, chanMsg, ms, "leave", config.CensorInvites) {
		return true, nil
	}

	return false, nil
}

// sendTemplate parses and executes the provided template, returns wether an error occured that we can retry from (temporary network failures and the like)
func sendTemplate(cs *dstate.ChannelState, tmpl string, ms *dstate.MemberState, name string, censorInvites bool) bool {
	ctx := templates.NewContext(cs.Guild, cs, ms)
	ctx.CurrentFrame.SendResponseInDM = cs.Type == discordgo.ChannelTypeDM

	ctx.Data["RealUsername"] = ms.Username
	if censorInvites {
		newUsername := common.ReplaceServerInvites(ms.Username, ms.Guild.ID, "[removed-server-invite]")
		if newUsername != ms.Username {
			ms.Username = newUsername + fmt.Sprintf("(user ID: %d)", ms.ID)
			ctx.Data["UsernameHasInvite"] = true
		}
	}

	msg, err := ctx.Execute(tmpl)

	if err != nil {
		logger.WithError(err).WithField("guild", cs.Guild.ID).Warnf("Failed parsing/executing %s template", name)
		return false
	}

	msg = strings.TrimSpace(msg)
	if msg == "" {
		return false
	}

	switch {
	case cs.Type == discordgo.ChannelTypeDM:
		_, err = common.BotSession.ChannelMessageSend(cs.ID, msg)
	case !ctx.CurrentFrame.DelResponse:
		send := ctx.MessageSend("")
		bot.QueueMergedMessage(cs.ID, msg, send.AllowedMentions)
	default:
		var m *discordgo.Message
		m, err = common.BotSession.ChannelMessageSendComplex(cs.ID, ctx.MessageSend(msg))
		if err == nil && ctx.CurrentFrame.DelResponse {
			templates.MaybeScheduledDeleteMessage(cs.Guild.ID, cs.ID, m.ID, ctx.CurrentFrame.DelResponseDelay)
		}
	}

	if err != nil {
		l := logger.WithError(err).WithField("guild", cs.Guild.ID)
		if common.IsDiscordErr(err, discordgo.ErrCodeCannotSendMessagesToThisUser) {
			l.Warn("Failed sending " + name)
		} else {
			l.Error("Failed sending " + name)
		}
	}

	return bot.CheckDiscordErrRetry(err)
}

func HandleChannelUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	cur := evt.ChannelUpdate()

	curChannel := bot.State.ChannelCopy(true, cur.ID)
	if curChannel == nil {
		return
	}

	oldTopic := curChannel.Topic

	logger.Infof("Old topic -> %#v || New topic -> %#v", oldTopic, cur.Topic)

	if oldTopic == cur.Topic {
		return
	}

	config, err := GetConfig(cur.GuildID)
	if err != nil {
		return true, errors.WithStackIf(err)
	}

	if !config.TopicEnabled {
		return
	}

	topicChannel := cur.Channel.ID
	if config.TopicChannelInt() != 0 {
		c := curChannel.Guild.Channel(true, config.TopicChannelInt())
		if c != nil {
			topicChannel = c.ID
		}
	}

	go analytics.RecordActiveUnit(cur.GuildID, &Plugin{}, "posted_topic_change")

	go func() {
		_, err := common.BotSession.ChannelMessageSend(topicChannel, fmt.Sprintf("Topic in channel <#%d> changed to **%s**", cur.ID, cur.Topic))
		if err != nil {
			logger.WithError(err).WithField("guild", cur.GuildID).Warn("Failed sending topic change message")
		}
	}()

	return false, nil
}
