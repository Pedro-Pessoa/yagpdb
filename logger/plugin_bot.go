package logger

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/analytics"
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/templates"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

var _ bot.BotInitHandler = (*Plugin)(nil)

// Initiates the events handler for discord logger
func (p *Plugin) BotInit() {
	// Channels
	eventsystem.AddHandlerAsyncLast(p, HandleChannelCreate, eventsystem.EventChannelCreate) // Channel creation
	eventsystem.AddHandlerAsyncLast(p, HandleChannelUpdate, eventsystem.EventChannelUpdate) // Channel update
	eventsystem.AddHandlerAsyncLast(p, HandleChannelDelete, eventsystem.EventChannelDelete) // Channel delete

	// Roles
	eventsystem.AddHandlerAsyncLast(p, HandleGuildRoleCreation, eventsystem.EventGuildRoleCreate) // Role creation
	eventsystem.AddHandlerAsyncLast(p, HandleGuildRoleUpdate, eventsystem.EventGuildRoleUpdate)   // Role update
	eventsystem.AddHandlerAsyncLast(p, HandleGuildRoleDelete, eventsystem.EventGuildRoleDelete)   // Role delete

	// Guilds
	eventsystem.AddHandlerAsyncLast(p, HandleGuildUpdate, eventsystem.EventGuildUpdate) // Guild update

	// Emojis
	eventsystem.AddHandlerAsyncLast(p, HandleGuildEmojisUpdate, eventsystem.EventGuildEmojisUpdate) // Emoji update

	// Messages
	eventsystem.AddHandlerAsyncLast(p, HandleMessageDelete, eventsystem.EventMessageDelete)         // Message delete
	eventsystem.AddHandlerAsyncLast(p, HandleMessageUpdate, eventsystem.EventMessageUpdate)         // Message edit
	eventsystem.AddHandlerAsyncLast(p, HandleMessageDeleteBulk, eventsystem.EventMessageDeleteBulk) // Message delete bulk

	// Members
	eventsystem.AddHandlerAsyncLast(p, HandleGuildMemberUpdate, eventsystem.EventGuildMemberUpdate) // Member update

	// Voice
	eventsystem.AddHandlerAsyncLast(p, HandleVoiceStateUpdate, eventsystem.EventVoiceStateUpdate) // Voice update
}

// Handles channels creations
func HandleChannelCreate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.ChannelCreate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.ChannelCreationEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) || logger.ServerLoggerMSG == "" {
		return false, nil
	}

	if len(logger.IgnoredLoggerCategories) > 0 && common.ContainsInt64Slice(logger.IgnoredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	if len(logger.RequiredLoggerCategories) > 0 && !common.ContainsInt64Slice(logger.RequiredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	cID := logger.DefaultLoggerChannelInt()
	if logger.ServerLoggerChannelInt() != 0 {
		cID = logger.ServerLoggerChannelInt()
	}

	senderChannel := evt.GS.Channel(true, cID)
	if senderChannel == nil {
		return false, nil
	}

	newChannel := dstate.NewChannelState(evt.GS, evt.GS, evt.ChannelCreate().Channel)

	go analytics.RecordActiveUnit(update.GuildID, &Plugin{}, "posted_channel_creation_msg")

	data := map[string]interface{}{
		"EventName":    "Channel Created",
		"EventChannel": templates.CtxChannelFromCSLocked(newChannel),
		"DefaultColor": 0x2effff,
		"EventType":    "ChannelEvent",
	}

	if sendTemplate(senderChannel, logger.ServerLoggerMSG, "channel_creation", data) {
		return true, nil
	}

	return false, nil
}

// Handles channels updates
func HandleChannelUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.ChannelUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.ChannelUpdatedEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) || logger.ServerLoggerMSG == "" {
		return false, nil
	}

	if len(logger.IgnoredLoggerCategories) > 0 && common.ContainsInt64Slice(logger.IgnoredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	if len(logger.RequiredLoggerCategories) > 0 && !common.ContainsInt64Slice(logger.RequiredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	cID := logger.DefaultLoggerChannelInt()
	if logger.ServerLoggerChannelInt() != 0 {
		cID = logger.ServerLoggerChannelInt()
	}

	bot.State.BeforeStateLocker.Lock()
	oldChannel := bot.State.BeforeStateChannelMap[update.ID]
	bot.State.BeforeStateLocker.Unlock()

	if oldChannel == nil || oldChannel.ChannelState == nil {
		return false, nil
	}

	valOld := reflect.ValueOf(oldChannel.ChannelState)
	typeOfOld := valOld.Type()

	valNew := reflect.ValueOf(update)
	typeOfNew := valNew.Type()

	changed := []interface{}{}

OUTER:
	for i := 0; i < valOld.NumField(); i++ {
		for j := 0; j < valNew.NumField(); j++ {
			if typeOfOld.Field(i).Name == typeOfNew.Field(j).Name {
				if valOld.Field(i).Interface() != valNew.Field(j).Interface() {
					changed = append(changed, typeOfOld.Field(i).Name)
					changed = append(changed, valOld.Field(i).Interface())
					changed = append(changed, valNew.Field(j).Interface())
				}

				continue OUTER
			}
		}
	}

	if len(changed) == 0 {
		return false, nil
	}

	senderChannel := evt.GS.Channel(true, cID)
	if senderChannel == nil {
		return false, nil
	}

	newChannel := dstate.NewChannelState(evt.GS, evt.GS, evt.ChannelCreate().Channel)

	go analytics.RecordActiveUnit(update.GuildID, &Plugin{}, "posted_channel_updated_msg")

	var changes string
	for i := 0; i < len(changed); i += 3 {
		changes = fmt.Sprintf("The channel **%s**'s changed from **%v** to **%v**\n", changed[i], changed[i+1], changed[i+2])
	}

	data := map[string]interface{}{
		"EventName":    "Channel Updated",
		"EventChannel": templates.CtxChannelFromCSLocked(newChannel),
		"DefaultColor": 0x8fffa5,
		"EventType":    "ChannelEvent",
		"Changes":      changes,
	}

	if sendTemplate(senderChannel, logger.ServerLoggerMSG, "channel_updated", data) {
		return true, nil
	}

	return false, nil
}

// Handles channels deletions
func HandleChannelDelete(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.ChannelDelete()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.ChannelDeletionEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) || logger.ServerLoggerMSG == "" {
		return false, nil
	}

	if len(logger.IgnoredLoggerCategories) > 0 && common.ContainsInt64Slice(logger.IgnoredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	if len(logger.RequiredLoggerCategories) > 0 && !common.ContainsInt64Slice(logger.RequiredLoggerCategories, update.Channel.ParentID) {
		return false, nil
	}

	cID := logger.DefaultLoggerChannelInt()
	if logger.ServerLoggerChannelInt() != 0 {
		cID = logger.ServerLoggerChannelInt()
	}

	senderChannel := evt.GS.Channel(true, cID)
	if senderChannel == nil {
		return false, nil
	}

	go analytics.RecordActiveUnit(update.GuildID, &Plugin{}, "posted_channel_deletion_msg")

	channel := templates.CtxChannel{
		Name:     update.Name,
		ParentID: update.ParentID,
	}

	data := map[string]interface{}{
		"EventName":    "Channel Deleted",
		"DefaultColor": 0xff6363,
		"EventType":    "ChannelEvent",
		"Channel":      channel,
	}

	if sendTemplate(senderChannel, logger.ServerLoggerMSG, "channel_deletion", data) {
		return true, nil
	}

	return false, nil
}

// Handles roles creations
func HandleGuildRoleCreation(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildRoleCreate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.RoleCreationEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) {
		return false, nil
	}

	cID := logger.DefaultLoggerChannelInt()
	if logger.ServerLoggerChannelInt() != 0 {
		cID = logger.ServerLoggerChannelInt()
	}

	senderChannel := evt.GS.Channel(true, cID)
	if senderChannel == nil {
		return false, nil
	}

	go analytics.RecordActiveUnit(update.GuildID, &Plugin{}, "posted_role_update_msg")

	data := map[string]interface{}{
		"EventName":    "Role Created",
		"DefaultColor": 0xbe27f5,
		"EventType":    "RoleEvent",
		"Role":         update.GuildRole,
	}

	if sendTemplate(senderChannel, logger.ServerLoggerMSG, "role_update", data) {
		return true, nil
	}

	return false, nil
}

// Handles roles updates
func HandleGuildRoleUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildRoleUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.RoleUpdateEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles roles deletions
func HandleGuildRoleDelete(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildRoleDelete()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.RoleDeletionEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles guild updates
func HandleGuildUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildUpdate()

	logger, err := GetLogger(update.ID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.ServerUpdateEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles emojis updates
func HandleGuildEmojisUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildEmojisUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.ServerEventsEnabled || !logger.EmojiUpdateEnbaled || (logger.DefaultLoggerChannelInt() == 0 && logger.ServerLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles messages deletes
func HandleMessageDelete(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.MessageDelete()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.MessageEventsEnabled || !logger.DeletedMessagesEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.MessageLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles messages updates/edits
func HandleMessageUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.MessageUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.MessageEventsEnabled || !logger.EditedMessagesEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.MessageLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles messages deletes bulk
func HandleMessageDeleteBulk(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.MessageDeleteBulk()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.MessageEventsEnabled || !logger.PurgedMessagesEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.MessageLoggerChannelInt() == 0) {
		return false, nil
	}

	return
}

// Handles members updates
func HandleGuildMemberUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.GuildMemberUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.MemberEventsEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.MemberLoggerChannelInt() == 0) || (!logger.MemberAvatarUpdateEnabled && !logger.MemberNameUpdateEnabled && logger.MemberRoleUpdateEnabled) {
		return false, nil
	}

	return
}

// Handles voice states updates
func HandleVoiceStateUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	if !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	update := evt.VoiceStateUpdate()

	logger, err := GetLogger(update.GuildID)
	if err != nil {
		return true, err
	}

	if !logger.LoggerEnabled || !logger.MemberEventsEnabled || (logger.DefaultLoggerChannelInt() == 0 && logger.VoiceLoggerChannelInt() == 0) || (!logger.JoinVoiceEnabled && !logger.LeaveVoiceEnabled && logger.SwapVoiceEnabled) {
		return false, nil
	}

	return
}

// sendTemplate parses and executes the provided template, returns wether an error occured that we can retry from (temporary network failures and the like)
func sendTemplate(cs *dstate.ChannelState, tmpl string, name string, data map[string]interface{}) bool {
	ctx := templates.NewContext(cs.Guild, cs, nil)
	ctx.CurrentFrame.SendResponseInDM = cs.Type == discordgo.ChannelTypeDM

	for k, v := range data {
		ctx.Data[k] = v
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

	_, err = common.BotSession.ChannelMessageSendComplex(cs.ID, ctx.MessageSend(msg))
	if err != nil {
		if common.IsDiscordErr(err, discordgo.ErrCodeCannotSendMessagesToThisUser) {
			logger.WithError(err).WithField("guild", cs.Guild.ID).Warn("Failed sending " + name)
		} else {
			logger.WithError(err).WithField("guild", cs.Guild.ID).Error("Failed sending " + name)
		}
	}

	return bot.CheckDiscordErrRetry(err)
}
