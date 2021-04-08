package templates

import (
	"time"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

// CtxChannel is almost a 1:1 copy of dstate.ChannelState, its needed because we cant axpose all those state methods
// we also cant use discordgo.Channel because that would likely break a lot of custom commands at this point.
type CtxChannel struct {
	// These fields never change
	ID        int64
	GuildID   int64
	IsPrivate bool

	Name                 string                           `json:"name"`
	Type                 discordgo.ChannelType            `json:"type"`
	Topic                string                           `json:"topic"`
	LastMessageID        int64                            `json:"last_message_id"`
	LastPinTimestamp     time.Time                        `json:"last_pin_timestamp"`
	NSFW                 bool                             `json:"nsfw"`
	Position             int                              `json:"position"`
	Bitrate              int                              `json:"bitrate"`
	PermissionOverwrites []*discordgo.PermissionOverwrite `json:"permission_overwrites"`
	ParentID             int64                            `json:"parent_id"`
	RateLimitPerUser     int                              `json:"rate_limit_per_user"`
	UserLimit            int                              `json:"user_limit"`
}

func CtxChannelFromCS(cs *dstate.ChannelState) *CtxChannel {
	ctxChannel := &CtxChannel{
		ID:                   cs.ID,
		IsPrivate:            cs.IsPrivate,
		Name:                 cs.Name,
		Type:                 cs.Type,
		Topic:                cs.Topic,
		LastMessageID:        cs.LastMessageID,
		LastPinTimestamp:     cs.LastPinTimestamp,
		NSFW:                 cs.NSFW,
		Position:             cs.Position,
		Bitrate:              cs.Bitrate,
		PermissionOverwrites: cs.PermissionOverwrites,
		ParentID:             cs.ParentID,
		RateLimitPerUser:     cs.RateLimitPerUser,
		UserLimit:            cs.UserLimit,
	}

	if !cs.IsPrivate {
		ctxChannel.GuildID = cs.Guild.ID
	}

	return ctxChannel
}

func CtxChannelFromCSLocked(cs *dstate.ChannelState) *CtxChannel {
	cs.Owner.RLock()
	defer cs.Owner.RUnlock()

	ctxChannel := &CtxChannel{
		ID:                   cs.ID,
		IsPrivate:            cs.IsPrivate,
		Name:                 cs.Name,
		Type:                 cs.Type,
		Topic:                cs.Topic,
		LastMessageID:        cs.LastMessageID,
		LastPinTimestamp:     cs.LastPinTimestamp,
		NSFW:                 cs.NSFW,
		Position:             cs.Position,
		Bitrate:              cs.Bitrate,
		PermissionOverwrites: cs.PermissionOverwrites,
		ParentID:             cs.ParentID,
		RateLimitPerUser:     cs.RateLimitPerUser,
		UserLimit:            cs.UserLimit,
	}

	if !cs.IsPrivate {
		ctxChannel.GuildID = cs.Guild.ID
	}

	return ctxChannel
}

func CtxChannelFromDGoChannel(dc *discordgo.Channel) *CtxChannel {
	ctxChannel := &CtxChannel{
		ID:                   dc.ID,
		Name:                 dc.Name,
		Type:                 dc.Type,
		Topic:                dc.Topic,
		LastMessageID:        dc.LastMessageID,
		NSFW:                 dc.NSFW,
		Position:             dc.Position,
		Bitrate:              dc.Bitrate,
		PermissionOverwrites: dc.PermissionOverwrites,
		ParentID:             dc.ParentID,
		RateLimitPerUser:     dc.RateLimitPerUser,
		UserLimit:            dc.UserLimit,
	}

	if !dstate.IsPrivate(ctxChannel.Type) {
		ctxChannel.GuildID = dc.GuildID
		ctxChannel.IsPrivate = false
	} else {
		ctxChannel.IsPrivate = true
	}

	ctxChannel.LastPinTimestamp, _ = dc.LastPinTimestamp.Parse()

	return ctxChannel
}

type CtxExecReturn struct {
	Return   []interface{}
	Response *discordgo.MessageSend
}

func (c CtxExecReturn) String() string {
	if c.Response != nil {
		return c.Response.Content
	}

	return ""
}

type CtxMember struct {
	GuildID      int64           `json:"guild_id"`
	JoinedAt     time.Time       `json:"joined_at"`
	Nick         string          `json:"nick"`
	Deaf         bool            `json:"deaf"`
	Mute         bool            `json:"mute"`
	User         *discordgo.User `json:"user"`
	Roles        []int64         `json:"roles"`
	PremiumSince time.Time       `json:"premium_since"`
	Pending      bool            `json:"pending"`

	//////////////////////////////////////
	// NON STANDARD MEMBER FIELDS BELOW //
	//////////////////////////////////////
	ID             int64                  `json:"id"`
	Username       string                 `json:"username"`
	Discriminator  int32                  `json:"discriminator"`
	AnimatedAvatar bool                   `json:"animated_avatar"`
	Bot            bool                   `json:"bot"`
	ClientStatus   discordgo.ClientStatus `json:"client_status"`

	//extra fields from Member State
	Status     dstate.PresenceStatus   `json:"status"`
	Activities []*dstate.LightActivity `json:"activity"`
}

func CtxMemberFromMS(ms *dstate.MemberState) *CtxMember {
	ctxMember := &CtxMember{
		GuildID:      ms.GuildID,
		JoinedAt:     ms.JoinedAt,
		Nick:         ms.Nick,
		Deaf:         ms.Deaf,
		Mute:         ms.Mute,
		User:         ms.DGoUser(),
		Roles:        ms.Roles,
		PremiumSince: ms.PremiumSince,
		Pending:      ms.Pending,

		ID:             ms.ID,
		Username:       ms.Username,
		Discriminator:  ms.Discriminator,
		AnimatedAvatar: ms.AnimatedAvatar,
		Bot:            ms.Bot,
		ClientStatus:   ms.ClientStatus,
		Status:         ms.PresenceStatus,
		Activities:     ms.PresenceActivities,
	}

	return ctxMember
}
