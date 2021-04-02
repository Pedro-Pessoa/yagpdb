package dstate

import (
	"time"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

// ChannelState represents a channel's state
type ChannelState struct {
	// These fields never change
	ID        int64       `json:"id"`
	Owner     RWLocker    `json:"-" msgpack:"-"`
	Guild     *GuildState `json:"-" msgpack:"-"`
	IsPrivate bool

	// Mutable fields, use Copy() or lock it
	Name                 string                           `json:"name"`
	Type                 discordgo.ChannelType            `json:"type"`
	Topic                string                           `json:"topic"`
	LastMessageID        int64                            `json:"last_message_id"`
	LastPinTimestamp     time.Time                        `json:"last_pin_timestamp"`
	NSFW                 bool                             `json:"nsfw"`
	Icon                 string                           `json:"icon"`
	Position             int                              `json:"position"`
	Bitrate              int                              `json:"bitrate"`
	PermissionOverwrites []*discordgo.PermissionOverwrite `json:"permission_overwrites"`
	ParentID             int64                            `json:"parent_id"`

	// Safe to access in a copy, but not write to in a copy
	Recipients []*discordgo.User `json:"recipients"`

	// Accessing the channel without locking the owner yields undefined behaviour
	Messages []*MessageState `json:"messages" msgpack:"-"`

	// The last message edit we didn't have the original message tracked for
	// we don't put those in the state because the ordering would be messed up
	// and there could be some unknown messages before and after
	// but in some cases (embed edits for example) the edit can come before the create event
	// for those edge cases we store the last edited unknown message here, then apply it as an update
	LastUnknownMsgEdit *discordgo.Message `json:"last_unknown_msg_edit"`

	// Amount of seconds a user has to wait before sending another message (0-21600)
	// bots, as well as users with the permission manage_messages or manage_channel, are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// ID of the DM creator Zeroed if guild channel
	OwnerID int64 `json:"owner_id"`

	// ApplicationID of the DM creator Zeroed if guild channel or not a bot user
	ApplicationID int64 `json:"application_id"`
}

func NewChannelState(guild *GuildState, owner RWLocker, channel *discordgo.Channel) *ChannelState {
	cs := &ChannelState{
		Owner: owner,
		Guild: guild,

		ID: channel.ID,
		// Type chan change, but the channel can never go from a dm type to a guild type, or vice versa
		// since its usefull to access this without locking, store that here
		IsPrivate: IsPrivate(channel.Type),

		Type:                 channel.Type,
		Name:                 channel.Name,
		Topic:                channel.Topic,
		LastMessageID:        channel.LastMessageID,
		NSFW:                 channel.NSFW,
		Position:             channel.Position,
		Bitrate:              channel.Bitrate,
		PermissionOverwrites: channel.PermissionOverwrites,
		ParentID:             channel.ParentID,
		OwnerID:              channel.OwnerID,
		ApplicationID:        channel.ApplicationID,
		RateLimitPerUser:     channel.RateLimitPerUser,

		Recipients: channel.Recipients,
	}

	cs.LastPinTimestamp, _ = channel.LastPinTimestamp.Parse()

	return cs
}

// DGoCopy returns a discordgo version of the channel representation
// usefull for legacy api's and whatnot
func (c *ChannelState) DGoCopy() *discordgo.Channel {
	channel := &discordgo.Channel{

		ID:   c.ID,
		Type: c.Type,

		Name:                 c.Name,
		Topic:                c.Topic,
		LastMessageID:        c.LastMessageID,
		NSFW:                 c.NSFW,
		Position:             c.Position,
		Bitrate:              c.Bitrate,
		PermissionOverwrites: c.PermissionOverwrites,
		ParentID:             c.ParentID,
		Recipients:           c.Recipients,
	}

	if c.Guild != nil {
		channel.GuildID = c.Guild.ID
	}

	return channel
}

// StrID is a conveniece method for retrieving the id in string form
func (cs *ChannelState) StrID() string {
	return discordgo.StrID(cs.ID)
}

// Recipient returns the channels recipient, if you modify this you get undefined behaviour
// This does no locking UNLESS this is a group dm
//
// In case of group dms, this will return the first recipient
func (cs *ChannelState) Recipient() *discordgo.User {
	if cs.Type == discordgo.ChannelTypeGroupDM {
		cs.Owner.RLock()
		defer cs.Owner.RUnlock()
	}
	if len(cs.Recipients) < 1 {
		return nil
	}

	return cs.Recipients[0]
}

// Copy returns a copy of the channel
// permissionoverwrites will be copied
// note: this is not a deep copy, modifying any of the slices is undefined behaviour,
// reading is fine as they're completely replaced when a update occurs
// (messages is another thing and is not available in this copy, manual management of the lock is needed for that)
func (c *ChannelState) Copy(lock bool) *ChannelState {
	if lock {
		c.Owner.RLock()
		defer c.Owner.RUnlock()
	}

	cop := new(ChannelState)
	*cop = *c

	cop.Messages = nil
	return cop
}

// Update updates a channel
// Undefined behaviour if owner (guild or state) is not locked
func (c *ChannelState) Update(lock bool, newChannel *discordgo.Channel) {
	if lock {
		c.Owner.Lock()
		defer c.Owner.Unlock()
	}

	if newChannel.PermissionOverwrites != nil {
		c.PermissionOverwrites = newChannel.PermissionOverwrites
	}

	if newChannel.Recipients != nil && c.Type == discordgo.ChannelTypeGroupDM {
		c.Recipients = newChannel.Recipients
	}

	c.Name = newChannel.Name
	c.Topic = newChannel.Topic
	c.LastMessageID = newChannel.LastMessageID
	c.NSFW = newChannel.NSFW
	c.Position = newChannel.Position
	c.Bitrate = newChannel.Bitrate
	c.ParentID = newChannel.ParentID
}

// Message returns a message REFERENCE by id or nil if none found
// The only field safe to query on a message reference without locking the owner (guild or state) is ID
func (c *ChannelState) Message(lock bool, mID int64) *MessageState {
	if lock {
		c.Owner.RLock()
		defer c.Owner.RUnlock()
	}

	index := c.messageIndex(mID)

	if index == -1 {
		return nil
	}

	return c.Messages[index]
}

// MessageCopy returns a copy of the message specified by id, its safe to read all fields, but it's not safe to modify any
func (c *ChannelState) MessageCopy(lock bool, mID int64) *MessageState {
	if lock {
		c.Owner.RLock()
		defer c.Owner.RUnlock()
	}

	index := c.messageIndex(mID)
	if index == -1 {
		return nil
	}

	return c.Messages[index].Copy()
}

func (c *ChannelState) messageIndex(mID int64) int {
	// since this should be ordered by low-high, maybe we should do a binary search?
	for i, v := range c.Messages {
		if v.ID == mID {
			return i
		}
	}

	return -1
}

// MessageAddUpdate adds or updates an existing message
func (c *ChannelState) MessageAddUpdate(lock bool, msg *discordgo.Message, edit bool) {
	if lock {
		c.Owner.Lock()
		defer c.Owner.Unlock()
	}

	existingIndex := c.messageIndex(msg.ID)

	if existingIndex != -1 {
		c.Messages[existingIndex].Update(msg)
		return
	}

	if edit {
		c.LastUnknownMsgEdit = msg
		return
	}

	ms := MessageStateFromMessage(msg)

	if c.LastUnknownMsgEdit != nil && c.LastUnknownMsgEdit.ID == ms.ID {
		ms.Update(c.LastUnknownMsgEdit)
		c.LastUnknownMsgEdit = nil
	}

	c.Messages = append(c.Messages, ms)
}

// MessageRemove removes a message from the channelstate
// If mark is true the the message will just be marked as deleted and not removed from the state
func (c *ChannelState) MessageRemove(lock bool, messageID int64, mark bool) {
	if lock {
		c.Owner.Lock()
		defer c.Owner.Unlock()
	}

	for i, ms := range c.Messages {
		if ms.ID == messageID {
			if !mark {
				c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			} else {
				ms.Deleted = true
			}
			return
		}
	}
}

const DiscordEpoch = 1420070400000

// assumes the owner is locked when ran
func (c *ChannelState) runGC(t time.Time, maxMessageAge time.Duration, maxMessages int) (messagesRemoved int) {
	c.Messages, messagesRemoved = clearMessageBuffer(c.Messages, t, maxMessageAge, maxMessages)
	return messagesRemoved
}

func clearMessageBuffer(buf []*MessageState, t time.Time, maxMessageAge time.Duration, maxMessages int) (newMessageBuffer []*MessageState, messagesRemoved int) {
	// check basic max message limit first
	if maxMessages > 0 && len(buf) > maxMessages {
		messagesRemoved += len(buf) - maxMessages
		buf = buf[len(buf)-maxMessages:]
	}

	if maxMessageAge < 1 || len(buf) < 1 {
		// nothing more to do....
		if len(buf) < 1 {
			buf = nil
		}
		return buf, messagesRemoved
	}

	// create a fake snowflake for comparing the message timestamps with fast
	fakeSnowflake := ((t.Add(-maxMessageAge).Unix() * 1000) - DiscordEpoch) << 22
	if buf[0].ID > fakeSnowflake {
		// fast path in case even the oldest message tracked is within max age, do nothing
		return buf, messagesRemoved
	}

	// Iterate reverse, new messages are at the end of the slice so iterate until we hit the first old message
	for i := len(buf) - 1; i >= 0; i-- {
		if buf[i].ID < fakeSnowflake {
			// older than the limit
			// all messages before this is old aswell
			messagesRemoved += i + 1
			buf = buf[i+1:]
			break
		}
	}

	if len(buf) < 1 {
		buf = nil
	}

	return buf, messagesRemoved
}
