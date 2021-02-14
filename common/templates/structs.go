package templates

import (
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
	NSFW                 bool                             `json:"nsfw"`
	Position             int                              `json:"position"`
	Bitrate              int                              `json:"bitrate"`
	PermissionOverwrites []*discordgo.PermissionOverwrite `json:"permission_overwrites"`
	ParentID             int64                            `json:"parent_id"`
}

func CtxChannelFromCS(cs *dstate.ChannelState) *CtxChannel {
	ctxChannel := &CtxChannel{
		ID:                   cs.ID,
		IsPrivate:            cs.IsPrivate,
		Name:                 cs.Name,
		Type:                 cs.Type,
		Topic:                cs.Topic,
		LastMessageID:        cs.LastMessageID,
		NSFW:                 cs.NSFW,
		Position:             cs.Position,
		Bitrate:              cs.Bitrate,
		PermissionOverwrites: cs.PermissionOverwrites,
		ParentID:             cs.ParentID,
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
		NSFW:                 cs.NSFW,
		Position:             cs.Position,
		Bitrate:              cs.Bitrate,
		PermissionOverwrites: cs.PermissionOverwrites,
		ParentID:             cs.ParentID,
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
	}

	if !dstate.IsPrivate(ctxChannel.Type) {
		ctxChannel.GuildID = dc.GuildID
		ctxChannel.IsPrivate = false
	} else {
		ctxChannel.IsPrivate = true
	}

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
	*discordgo.Member
	//extra fields from Member State
	Status     dstate.PresenceStatus   `json:"status"`
	Activities []*dstate.LightActivity `json:"activity"`
}

func CtxMemberFromMS(ms *dstate.MemberState) *CtxMember {
	ctxMember := &CtxMember{
		Member:     ms.DGoCopy(),
		Status:     ms.PresenceStatus,
		Activities: ms.PresenceActivities,
	}

	return ctxMember
}
