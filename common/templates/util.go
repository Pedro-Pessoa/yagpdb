package templates

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
)

func emojiArg(emoji string) string {
	emoji = strings.TrimPrefix(emoji, "<")
	emoji = strings.TrimSuffix(emoji, ">")
	emoji = strings.TrimPrefix(emoji, ":")

	return emoji
}

func channelValueValidation(input interface{}) (out string, valid bool) {
	switch t := input.(type) {
	case string:
		_, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return "", false
		}

		return t, true
	default:
		newValue := ToString(t)
		if newValue == "" {
			return "", false
		}

		_, err := strconv.ParseInt(newValue, 10, 64)
		if err != nil {
			return "", false
		}

		return newValue, true
	}
}

var permMap = map[string]int64{
	"view_channel":           discordgo.PermissionViewChannel,
	"read_messages":          discordgo.PermissionViewChannel, // deprecated, let here for convinience
	"send_messages":          discordgo.PermissionSendMessages,
	"manage_messages":        discordgo.PermissionManageMessages,
	"send_tts_messages":      discordgo.PermissionSendTTSMessages,
	"embed_links":            discordgo.PermissionEmbedLinks,
	"attach_files":           discordgo.PermissionAttachFiles,
	"read_message_history":   discordgo.PermissionReadMessageHistory,
	"mention_everyone":       discordgo.PermissionMentionEveryone,
	"use_external_emojis":    discordgo.PermissionUseExternalEmojis,
	"voice_connect":          discordgo.PermissionVoiceConnect,
	"voice_speak":            discordgo.PermissionVoiceSpeak,
	"voice_mute_members":     discordgo.PermissionVoiceMuteMembers,
	"voice_deafen_members":   discordgo.PermissionVoiceDeafenMembers,
	"voice_move_members":     discordgo.PermissionVoiceMoveMembers,
	"voice_use_vad":          discordgo.PermissionVoiceUseVAD,
	"voice_priority_speaker": discordgo.PermissionVoicePrioritySpeaker,
	"add_reactions":          discordgo.PermissionAddReactions,
	"create_invite":          discordgo.PermissionCreateInstantInvite,
	"manage_channel":         discordgo.PermissionManageChannels,
	"manage_webhooks":        discordgo.PermissionManageWebhooks,
}

func getLen(from interface{}) int {
	v := reflect.ValueOf(from)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len()
	default:
		return 0
	}
}

func targetUserID(input interface{}) int64 {
	switch t := input.(type) {
	case *discordgo.User:
		return t.ID
	case string:
		str := strings.TrimSpace(t)
		if strings.HasPrefix(str, "<@") && strings.HasSuffix(str, ">") && (len(str) > 4) {
			trimmed := str[2 : len(str)-1]
			if trimmed[0] == '!' {
				trimmed = trimmed[1:]
			}
			str = trimmed
		}

		return ToInt64(str)
	default:
		return ToInt64(input)
	}
}

func (c *Context) getMember(id interface{}) (*dstate.MemberState, error) {
	targetID := targetUserID(id)
	if targetID == 0 {
		return nil, fmt.Errorf("Target %v not found", id)
	}

	ms, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return nil, err
	}

	if ms == nil {
		return nil, errors.New("MemberState not found")
	}

	return ms, nil
}

// ChannelArg converts a variety of types of argument into a channel, verifying that it exists
func (c *Context) ChannelArg(v interface{}) int64 {
	c.GS.RLock()
	defer c.GS.RUnlock()

	// Look for the channel
	if v == nil && c.CurrentFrame.CS != nil {
		// No channel passed, assume current channel
		return c.CurrentFrame.CS.ID
	}

	return c.VerifyChannel(v)
}

// ChannelArgNoDM is the same as ChannelArg but will not accept DM channels
func (c *Context) ChannelArgNoDM(v interface{}) int64 {
	c.GS.RLock()
	defer c.GS.RUnlock()

	// Look for the channel
	if v == nil && c.CurrentFrame.CS != nil {
		// No channel passed, assume current channel
		v = c.CurrentFrame.CS.ID
	}

	return c.VerifyChannel(v)
}

func (c *Context) VerifyChannel(v interface{}) int64 {
	var cid int64
	switch t := v.(type) {
	case int, int64:
		// Channel id passed
		cid = ToInt64(t)
	case string:
		parsed, err := strconv.ParseInt(t, 10, 64)
		if err == nil {
			// Channel id passed in string format
			cid = parsed
		} else {
			// See if it is a channel mention
			if strings.HasPrefix(t, "<#") && strings.HasSuffix(t, ">") && (len(t) > 3) {
				parsedMention, err := strconv.ParseInt(t[2:len(t)-1], 10, 64)
				if err == nil {
					cid = parsedMention
				}
			} else { // Channel name, look for it
				for _, v := range c.GS.Channels {
					if strings.EqualFold(t, v.Name) && v.Type == discordgo.ChannelTypeGuildText {
						return v.ID // Channel found and is part of the Guild, we can return already
					}
				}
				return 0 // If we got here it means the channel provided was a name and it is not part of the guild, so we can return.
			}
		}
	default:
		return 0 // Invalid channel provided
	}

	// Make sure the channel is part of the guild
	for k := range c.GS.Channels {
		if k == cid {
			return cid
		}
	}

	return 0
}

func (c *Context) checkSafeStringDictNoRecursion(d SDict, n int) bool {
	if n > 1000 {
		return false
	}

	for _, v := range d {
		if cast, ok := v.(Dict); ok {
			if !c.checkSafeDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if cast, ok := v.(SDict); ok {
			if !c.checkSafeStringDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if reflect.DeepEqual(v, c.Data) {
			return false
		}
	}

	return true
}

func (c *Context) checkSafeDictNoRecursion(d Dict, n int) bool {
	if n > 1000 {
		return false
	}

	for _, v := range d {
		if cast, ok := v.(Dict); ok {
			if !c.checkSafeDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if cast, ok := v.(SDict); ok {
			if !c.checkSafeStringDictNoRecursion(cast, n+1) {
				return false
			}
		}

		if reflect.DeepEqual(v, c.Data) {
			return false
		}
	}

	return true
}
