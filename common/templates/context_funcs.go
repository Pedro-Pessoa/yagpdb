package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/scheduledevents2"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/pkgs/template"
)

var (
	ErrTooManyCalls    = errors.New("Too many calls to this function")
	ErrTooManyAPICalls = errors.New("Too many potential discord api calls function")
	reflectValueType   = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

// Message Functions

func (c *Context) buildDM(gName string, s ...interface{}) *discordgo.MessageSend {
	info := "DM enviada pelo servidor **" + gName + "**"

	msgSend := &discordgo.MessageSend{
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	}

	switch t := s[0].(type) {
	case *discordgo.MessageEmbed:
		msgSend.Embed = t
	case *discordgo.MessageSend:
		msgSend = t
		if (strings.TrimSpace(msgSend.Content) == "") && (msgSend.File == nil) {
			return nil
		}
	default:
		msgSend.Content = fmt.Sprint(s...)
	}

	if !bot.IsGuildWhiteListed(c.GS.ID) {
		if msgSend.Embed != nil {
			msgSend.Embed.Footer = &discordgo.MessageEmbedFooter{
				Text: info,
			}
		} else {
			msgSend.Content = info + "\n" + msgSend.Content
		}
	}

	return msgSend
}

func (c *Context) tmplSendDM(s ...interface{}) string {
	if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 1) || c.IncreaseCheckGenericAPICall() || c.MS == nil || c.IsExecedByLeaveMessage {
		return ""
	}

	c.GS.RLock()
	memberID, gName := c.MS.ID, c.GS.Guild.Name
	c.GS.RUnlock()

	msgSend := c.buildDM(gName, s...)
	if msgSend == nil {
		return ""
	}

	_, _ = bot.SendDMComplex(memberID, msgSend)
	return ""
}

func (c *Context) tmplSendDMWithError(s ...interface{}) (string, error) {
	if len(s) < 1 {
		return "", errors.New("No argument provided to SendDM")
	}

	if c.IncreaseCheckCallCounter("send_dm", 1) {
		return "", ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IsExecedByLeaveMessage {
		return "", errors.New("Can't use sendDM on leave msg")
	}

	if c.MS == nil {
		return "", errors.New("SendDM called on context with nil MemberState")
	}

	c.GS.RLock()
	memberID, gName := c.MS.ID, c.GS.Guild.Name
	c.GS.RUnlock()

	msgSend := c.buildDM(gName, s...)
	if msgSend == nil {
		return "", errors.New("Failed building the DM, join the support server if you need help")
	}

	_, err := bot.SendDMComplex(memberID, msgSend)
	return "", err
}

func (c *Context) tmplSendTargetDM(target interface{}, s ...interface{}) string {
	if bot.IsGuildWhiteListed(c.GS.ID) {
		if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 1) || c.IncreaseCheckGenericAPICall() || c.IsExecedByLeaveMessage {
			return ""
		}

		ts, _ := c.getMember(target)
		if ts == nil {
			return ""
		}

		msgSend := c.buildDM("", s...)
		if msgSend == nil {
			return ""
		}

		_, _ = bot.SendDMComplex(ts.ID, msgSend)
	}

	return ""
}

func (c *Context) tmplSendTargetDMWithError(target interface{}, s ...interface{}) (string, error) {
	if bot.IsGuildWhiteListed(c.GS.ID) {
		if len(s) < 1 {
			return "", errors.New("No argument provided to SendDM")
		}

		if c.IncreaseCheckCallCounter("send_dm", 1) {
			return "", ErrTooManyCalls
		}

		if c.IncreaseCheckGenericAPICall() {
			return "", ErrTooManyAPICalls
		}

		if c.IsExecedByLeaveMessage {
			return "", errors.New("Can't use sendTargetDM on leave msg")
		}

		ts, err := c.getMember(target)
		if err != nil {
			return "", err
		}

		msgSend := c.buildDM("", s...)
		if msgSend == nil {
			return "", errors.New("Failed building the DM, join the support server if you need help")
		}

		_, err = bot.SendDMComplex(ts.ID, msgSend)
		return "", err
	}

	return "", nil
}

func (c *Context) tmplSendMessage(filterSpecialMentions bool, returnID bool, ignoreError bool) func(channel interface{}, msg interface{}, replyData ...interface{}) (interface{}, error) {
	parseMentions := []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers}
	if !filterSpecialMentions {
		parseMentions = append(parseMentions, discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeEveryone)
	}

	return func(channel interface{}, msg interface{}, replyData ...interface{}) (interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("send_msg", 10, 20) {
			if ignoreError {
				return "", nil
			}

			return nil, ErrTooManyCalls
		}

		if c.IncreaseCheckGenericAPICall() {
			if ignoreError {
				return "", nil
			}

			return "", ErrTooManyAPICalls
		}

		cid := c.ChannelArg(channel)
		if cid == 0 {
			if ignoreError {
				return "", nil
			}

			return "", errors.New("Invalid channel provided for SendMessage")
		}

		isDM := cid != c.ChannelArgNoDM(channel)
		c.GS.RLock()
		info := "DM enviada pelo servidor **" + c.GS.Guild.Name + "**"
		c.GS.RUnlock()
		WL := bot.IsGuildWhiteListed(c.GS.ID)

		var m *discordgo.Message
		msgSend := &discordgo.MessageSend{
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: parseMentions,
			},
		}

		var err error

		switch typedMsg := msg.(type) {
		case *discordgo.MessageEmbed:
			if isDM && !WL {
				typedMsg.Footer = &discordgo.MessageEmbedFooter{
					Text: info,
				}
			}
			msgSend.Embed = typedMsg
		case *discordgo.MessageSend:
			msgSend = typedMsg

			if !filterSpecialMentions {
				msgSend.AllowedMentions = &discordgo.MessageAllowedMentions{Parse: parseMentions}
			}

			if isDM && !WL {
				if typedMsg.Embed != nil {
					typedMsg.Embed.Footer = &discordgo.MessageEmbedFooter{
						Text: info,
					}
				} else {
					typedMsg.Content = info + "\n" + typedMsg.Content
				}
			}
		default:
			if isDM && !WL {
				msgSend.Content = info + "\n" + ToString(msg)
			} else {
				msgSend.Content = ToString(msg)
			}
		}

		var dict SDict
		switch len(replyData) {
		case 0:
			m, err = common.BotSession.ChannelMessageSendComplex(cid, msgSend)
			if err != nil {
				if ignoreError {
					return "", nil
				}

				return "", err
			}

			if returnID {
				return m.ID, nil
			}

			return "", nil
		case 1:
			val := reflect.ValueOf(replyData[0])
			switch val.Kind() {
			case reflect.Map:
				dict, err = StringKeyDictionary(replyData[0])
				if err != nil {
					if ignoreError {
						return "", nil
					}

					return "", err
				}
			default:
				if ignoreError {
					return "", nil
				}

				return "", errors.Errorf("Invalid argument for ReplyData of type %s. Must be a map", val.Type().Name())
			}
		default:
			dict, err = StringKeyDictionary(replyData...)
			if err != nil {
				if ignoreError {
					return "", nil
				}

				return "", err
			}
		}

		replyChannelID := c.ChannelArg(dict.Get("channel_id"))
		replyMessageID := ToInt64(dict.Get("message_id"))

		if replyChannelID == 0 || replyMessageID == 0 {
			if ignoreError {
				return "", nil
			}

			return "", errors.New("Invalid channel or message ID provided for ReplyData")
		}

		reference := &discordgo.MessageReference{
			ChannelID: replyChannelID,
			MessageID: replyMessageID,
		}

		if !isDM {
			reference.GuildID = c.GS.ID
		}

		msgSend.Reference = reference

		m, err = common.BotSession.ChannelMessageSendComplex(cid, msgSend)
		if err != nil {
			if ignoreError {
				return "", nil
			}

			return "", err
		}

		if returnID {
			return m.ID, nil
		}

		return "", nil
	}
}

func (c *Context) tmplEditMessage(filterSpecialMentions bool) func(channel interface{}, msgID interface{}, msg interface{}) (interface{}, error) {
	return func(channel interface{}, msgID interface{}, msg interface{}) (interface{}, error) {
		if c.IncreaseCheckGenericAPICall() {
			return "", ErrTooManyAPICalls
		}

		cid := c.ChannelArgNoDM(channel)
		if cid == 0 {
			return "", errors.New("Unknown channel")
		}

		mID := ToInt64(msgID)
		msgEdit := &discordgo.MessageEdit{
			ID:      mID,
			Channel: cid,
		}

		var err error

		switch typedMsg := msg.(type) {
		case *discordgo.MessageEmbed:
			msgEdit.Embed = typedMsg
		case *discordgo.MessageEdit:
			//If both Embed and string are explicitly set as null, give an error message.
			if typedMsg.Content != nil && strings.TrimSpace(*typedMsg.Content) == "" && typedMsg.Embed != nil && typedMsg.Embed.GetMarshalNil() {
				return "", errors.New("both content and embed cannot be null")
			}

			msgEdit.Content = typedMsg.Content
			msgEdit.Embed = typedMsg.Embed
			msgEdit.AllowedMentions = typedMsg.AllowedMentions
		default:
			temp := ToString(msg)
			msgEdit.Content = &temp
		}

		if !filterSpecialMentions {
			msgEdit.AllowedMentions = &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers, discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeEveryone},
			}
		}

		_, err = common.BotSession.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			return "", err
		}

		return "", nil
	}
}

func (c *Context) tmplDelResponse(args ...interface{}) string {
	dur := 10
	if len(args) > 0 {
		dur = ToInt(args[0])
	}

	if dur > 86400 {
		dur = 86400
	}

	c.CurrentFrame.DelResponseDelay = dur
	c.CurrentFrame.DelResponse = true

	return ""
}

func (c *Context) tmplDelTrigger(args ...interface{}) string {
	if c.Msg != nil {
		return c.tmplDelMessage(c.Msg.ChannelID, c.Msg.ID, args...)
	}

	return ""
}

func (c *Context) tmplDelMessage(channel, msgID interface{}, args ...interface{}) string {
	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return ""
	}

	mID := ToInt64(msgID)

	dur := 10
	if len(args) > 0 {
		dur = ToInt(args[0])
	}

	if dur > 86400 {
		dur = 86400
	}

	MaybeScheduledDeleteMessage(c.GS.ID, cID, mID, dur)

	return ""
}

func (c *Context) tmplGetMessage(channel, msgID interface{}) (*discordgo.Message, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return nil, nil
	}

	mID := ToInt64(msgID)

	return common.BotSession.ChannelMessage(cID, mID)
}

func (c *Context) tmplGetMessageReactors(channel, msg interface{}, emoji string, limit int, before, after interface{}) ([]*discordgo.User, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	if c.IncreaseCheckCallCounterPremium("msg_reactors", 2, 5) {
		return nil, ErrTooManyCalls
	}

	channelID := c.ChannelArgNoDM(channel)
	if channelID == 0 {
		return nil, errors.New("Channel not found")
	}

	msgID := ToInt64(msg)
	if msgID == 0 {
		return nil, errors.New("Msg ID is invalid")
	}

	emojiID := emojiArg(emoji)
	if emojiID == "" {
		return nil, errors.New("Invalid emoji name provided")
	}

	switch {
	case limit > 100:
		return nil, errors.New("Limit can not be bigger than 100")
	case limit < 0:
		return nil, errors.New("Limit can not be negative")
	}

	beforeID := ToInt64(before)
	afterID := ToInt64(after)

	reactors, err := common.BotSession.MessageReactions(channelID, msgID, emojiID, limit, beforeID, afterID)
	if err != nil {
		return nil, err
	}

	return reactors, nil
}

func (c *Context) tmplPinMessage(channel, message interface{}) (interface{}, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckCallCounterPremium("msg_pin", 2, 5) {
		return nil, ErrTooManyCalls
	}

	mID := ToInt64(message)
	if mID == 0 {
		return "", errors.New("Invalid message provided")
	}

	channelID := c.ChannelArgNoDM(channel)
	if channelID == 0 {
		return "", errors.New("Invalid channel provided")
	}

	err := common.BotSession.ChannelMessagePin(channelID, mID)
	return "", err
}

func (c *Context) tmplUnpinMessage(channel, message interface{}) (interface{}, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckCallCounterPremium("msg_pin", 2, 5) {
		return nil, ErrTooManyCalls
	}

	mID := ToInt64(message)
	if mID == 0 {
		return "", errors.New("Invalid message provided")
	}

	channelID := c.ChannelArgNoDM(channel)
	if channelID == 0 {
		return "", errors.New("Invalid channel provided")
	}

	err := common.BotSession.ChannelMessageUnpin(channelID, mID)
	return "", err
}

// Templates functions

func (c *Context) tmplSendTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(channel, false, false, name, data...)
}

func (c *Context) tmplSendTemplateDM(name string, data ...interface{}) (interface{}, error) {
	if c.IsExecedByLeaveMessage {
		return "", errors.New("Can't use sendDM on leave msg")
	}

	return c.sendNestedTemplate(nil, true, false, name, data...)
}

func (c *Context) tmplWalkTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(channel, false, true, name, data...)
}

func (c *Context) tmplAddReturn(data ...interface{}) (interface{}, error) {
	if !c.CurrentFrame.isNestedTemplate || !c.CurrentFrame.execMode {
		return "", errors.New("Can only be used in nested templates in exec mode.")
	}

	if len(c.CurrentFrame.execReturn)+len(data) > 10 {
		return "", errors.New("Return length cannot exceed 10")
	}

	c.CurrentFrame.execReturn = append(c.CurrentFrame.execReturn, data...)
	return "", nil
}

func (c *Context) sendNestedTemplate(channel interface{}, dm, exec bool, name string, data ...interface{}) (interface{}, error) {
	if c.IncreaseCheckCallCounter("exec_child", 3) {
		return "", ErrTooManyCalls
	}

	if name == "" {
		return "", errors.New("No template name passed")
	}

	if c.CurrentFrame.isNestedTemplate {
		return "", errors.New("Can't call this in a nested template")
	}

	t := c.CurrentFrame.parsedTemplate.Lookup(name)
	if t == nil {
		return "", errors.New("Unknown template")
	}

	var cs *dstate.ChannelState
	// find the new context channel
	if !dm {
		if channel == nil {
			cs = c.CurrentFrame.CS
		} else {
			cID := c.ChannelArg(channel)
			if cID == 0 {
				return "", errors.New("Unknown channel")
			}

			cs = c.GS.ChannelCopy(true, cID)
			if cs == nil {
				return "", errors.New("Unknown channel")
			}
		}
	} else {
		if c.CurrentFrame.SendResponseInDM {
			cs = c.CurrentFrame.CS
		} else {
			ch, err := common.BotSession.UserChannelCreate(c.MS.ID)
			if err != nil {
				return "", err
			}

			cs = &dstate.ChannelState{
				Owner: c.GS,
				Guild: c.GS,
				ID:    ch.ID,
				Name:  c.MS.Username,
				Type:  discordgo.ChannelTypeDM,
			}
		}
	}

	oldFrame := c.newContextFrame(cs)
	defer func() {
		c.CurrentFrame = oldFrame
	}()

	if dm {
		c.CurrentFrame.SendResponseInDM = oldFrame.SendResponseInDM
	} else if channel == nil {
		// inherit
		c.CurrentFrame.SendResponseInDM = oldFrame.SendResponseInDM
	}

	c.CurrentFrame.execMode = exec
	// pass some data
	if len(data) > 1 {
		dict, _ := Dictionary(data...)
		c.Data["TemplateArgs"] = dict
		if !c.checkSafeDictNoRecursion(dict, 0) {
			return nil, errors.New("trying to pass the entire current context data in as templateargs, this is not needed, just use nil and access all other data normally")
		}
	} else if len(data) == 1 {
		if cast, ok := data[0].(map[string]interface{}); ok && reflect.DeepEqual(cast, c.Data) {
			return nil, errors.New("trying to pass the entire current context data in as templateargs, this is not needed, just use nil and access all other data normally")
		}
		c.Data["TemplateArgs"] = data[0]
	}

	// and finally execute the child template
	c.CurrentFrame.parsedTemplate = t
	resp, err := c.executeParsed()
	if err != nil {
		return "", err
	}

	if exec {
		var execReturnStruct CtxExecReturn
		execReturnStruct.Response = c.MessageSend(resp)
		execReturnStruct.Return = c.CurrentFrame.execReturn
		return execReturnStruct, err
	}

	m, err := c.SendResponse(resp)
	if err != nil {
		return "", err
	}

	if m != nil {
		return m.ID, err
	}

	return "", err
}

// Mention Functions

func (c *Context) tmplMentionEveryone() string {
	c.CurrentFrame.MentionEveryone = true
	return "@everyone"
}

func (c *Context) tmplMentionHere() string {
	c.CurrentFrame.MentionHere = true
	return "@here"
}

// Roles functions

// c.FindRole accepts all possible role inputs (names, IDs and mentions)
// and tries to find them on the current context
func (c *Context) FindRole(role interface{}) *discordgo.Role {
	switch t := role.(type) {
	case string:
		parsed, err := strconv.ParseInt(t, 10, 64)
		if err == nil {
			return c.GS.RoleCopy(true, ToInt64(parsed))
		}

		if strings.HasPrefix(t, "<@&") && strings.HasSuffix(t, ">") && (len(t) > 4) {
			parsedMention, err := strconv.ParseInt(t[2:len(t)-1], 10, 64)
			if err == nil {
				return c.GS.RoleCopy(true, ToInt64(parsedMention))
			}
		}

		return c.GS.RoleCopyByName(true, t)
	default:
		int64Role := ToInt64(t)
		if int64Role == 0 {
			return nil
		}

		return c.GS.RoleCopy(true, int64Role)
	}
}

func (c *Context) tmplGetRole(r interface{}) (*discordgo.Role, error) {
	if c.IncreaseCheckStateLock() {
		return nil, ErrTooManyCalls
	}

	return c.FindRole(r), nil
}

func (c *Context) tmplMentionRole(roleInput interface{}) string {
	if c.IncreaseCheckStateLock() {
		return "(too many state locks)"
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return "(role not found)"
	}

	if common.ContainsInt64Slice(c.CurrentFrame.MentionRoles, role.ID) {
		return role.Mention()
	}

	c.CurrentFrame.MentionRoles = append(c.CurrentFrame.MentionRoles, role.ID)
	return role.Mention()
}

func (c *Context) tmplMentionRoleID(roleID interface{}) string {
	return c.tmplMentionRole(roleID)
}

func (c *Context) tmplMentionRoleName(roleName string) string {
	return c.tmplMentionRole(roleName)
}

func (c *Context) hasRole(roleInput interface{}) (bool, error) {
	if c.IncreaseCheckStateLock() {
		return false, ErrTooManyCalls
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return false, fmt.Errorf("Role %v not found", roleInput)
	}

	return common.ContainsInt64Slice(c.MS.Roles, role.ID), nil
}

func (c *Context) tmplHasRole(roleInput interface{}) bool {
	out, _ := c.hasRole(roleInput)
	return out
}

func (c *Context) tmplHasRoleID(roleID interface{}) bool {
	return c.tmplHasRole(roleID)
}

func (c *Context) tmplHasRoleName(roleName string) bool {
	return c.tmplHasRole(roleName)
}

func (c *Context) tmplHasRoleWithError(roleInput interface{}) (bool, error) {
	return c.hasRole(roleInput)
}

func (c *Context) tmplHasRoleIDWithError(roleID interface{}) (bool, error) {
	return c.tmplHasRoleWithError(roleID)
}

func (c *Context) tmplHasRoleNameWithError(roleName string) (bool, error) {
	return c.tmplHasRoleWithError(roleName)
}

func (c *Context) targetHasRole(target interface{}, roleInput interface{}) (bool, error) {
	if c.IncreaseCheckStateLock() {
		return false, ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return false, ErrTooManyAPICalls
	}

	ts, err := c.getMember(target)
	if err != nil {
		return false, err
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return false, fmt.Errorf("Role %v not found", roleInput)
	}

	return common.ContainsInt64Slice(ts.Roles, role.ID), nil
}

func (c *Context) tmplTargetHasRole(target interface{}, roleInput interface{}) bool {
	out, _ := c.targetHasRole(target, roleInput)
	return out
}

func (c *Context) tmplTargetHasRoleID(target interface{}, roleID interface{}) bool {
	return c.tmplTargetHasRole(target, roleID)

}

func (c *Context) tmplTargetHasRoleName(target interface{}, roleName string) bool {
	return c.tmplTargetHasRole(target, roleName)
}

func (c *Context) tmplTargetHasRoleWithError(target interface{}, roleInput interface{}) (bool, error) {
	return c.targetHasRole(target, roleInput)
}

func (c *Context) tmplTargetHasRoleIDWithError(target interface{}, roleID interface{}) (bool, error) {
	return c.tmplTargetHasRoleWithError(target, roleID)

}

func (c *Context) tmplTargetHasRoleNameWithError(target interface{}, roleName string) (bool, error) {
	return c.tmplTargetHasRoleWithError(target, roleName)
}

func (c *Context) giveRole(target interface{}, roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckStateLock() {
		return "", ErrTooManyCalls
	}

	var delay int
	if len(optionalArgs) > 0 {
		delay = ToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return "", fmt.Errorf("Target %v not found", target)
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return "", fmt.Errorf("Role %v not found", roleInput)
	}

	if delay > 0 {
		err := scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
		if err != nil {
			return "", err
		}
	} else {
		ms, err := bot.GetMember(c.GS.ID, targetID)
		var hasRole bool
		if ms != nil && err == nil {
			hasRole = common.ContainsInt64Slice(ms.Roles, role.ID)
		}

		if hasRole {
			// User already has this role, nothing to be done
			return "", nil
		}

		err = common.BotSession.GuildMemberRoleAdd(c.GS.ID, targetID, role.ID)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplGiveRole(target interface{}, roleInput interface{}, optionalArgs ...interface{}) string {
	out, _ := c.giveRole(target, roleInput, optionalArgs...)
	return out
}

func (c *Context) tmplGiveRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	return c.tmplGiveRole(target, roleID, optionalArgs...)
}

func (c *Context) tmplGiveRoleName(target interface{}, roleName string, optionalArgs ...interface{}) string {
	return c.tmplGiveRole(target, roleName, optionalArgs...)
}

func (c *Context) tmplGiveRoleWithError(target interface{}, roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	return c.giveRole(target, roleInput, optionalArgs...)
}

func (c *Context) tmplGiveRoleIDWithError(target interface{}, roleID interface{}, optionalArgs ...interface{}) (string, error) {
	return c.tmplGiveRoleWithError(target, roleID, optionalArgs...)
}

func (c *Context) tmplGiveRoleNameWithError(target interface{}, roleName string, optionalArgs ...interface{}) (string, error) {
	return c.tmplGiveRoleWithError(target, roleName, optionalArgs...)
}

func (c *Context) addRole(roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckStateLock() {
		return "", ErrTooManyCalls
	}

	var delay int
	if len(optionalArgs) > 0 {
		delay = ToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", errors.New("tmplAddRole called on context with nil MemberState")
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return "", fmt.Errorf("Role %v not found", roleInput)
	}

	if delay > 0 {
		err := scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
		if err != nil {
			return "", err
		}
	} else {
		err := common.AddRoleDS(c.MS, role.ID)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplAddRole(roleInput interface{}, optionalArgs ...interface{}) string {
	out, _ := c.addRole(roleInput, optionalArgs...)
	return out
}

func (c *Context) tmplAddRoleID(roleID interface{}, optionalArgs ...interface{}) string {
	return c.tmplAddRole(roleID, optionalArgs...)
}

func (c *Context) tmplAddRoleName(roleName string, optionalArgs ...interface{}) string {
	return c.tmplAddRole(roleName, optionalArgs...)
}

func (c *Context) tmplAddRoleWithError(roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	return c.addRole(roleInput, optionalArgs...)
}

func (c *Context) tmplAddRoleIDWithError(roleID interface{}, optionalArgs ...interface{}) (string, error) {
	return c.tmplAddRoleWithError(roleID, optionalArgs...)
}

func (c *Context) tmplAddRoleNameWithError(roleName string, optionalArgs ...interface{}) (string, error) {
	return c.tmplAddRoleWithError(roleName, optionalArgs...)
}

func (c *Context) takeRole(target interface{}, roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckStateLock() {
		return "", ErrTooManyCalls
	}

	var delay int
	if len(optionalArgs) > 0 {
		delay = ToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return "", fmt.Errorf("Target %v not found", target)
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return "", fmt.Errorf("Role %v not found", roleInput)
	}

	if delay > 0 {
		err := scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
		if err != nil {
			return "", err
		}
	} else {
		ms, err := bot.GetMember(c.GS.ID, targetID)
		hasRole := true
		if ms != nil && err == nil {
			hasRole = common.ContainsInt64Slice(ms.Roles, role.ID)
		}

		if !hasRole {
			// User does not have the role, nothing to be done
			return "", nil
		}

		err = common.BotSession.GuildMemberRoleRemove(c.GS.ID, targetID, role.ID)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplTakeRole(target interface{}, roleInput interface{}, optionalArgs ...interface{}) string {
	out, _ := c.takeRole(target, roleInput, optionalArgs...)
	return out
}

func (c *Context) tmplTakeRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	return c.tmplTakeRole(target, roleID, optionalArgs...)
}

func (c *Context) tmplTakeRoleName(target interface{}, roleName string, optionalArgs ...interface{}) string {
	return c.tmplTakeRole(target, roleName, optionalArgs...)
}

func (c *Context) tmplTakeRoleWithError(target interface{}, roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	return c.takeRole(target, roleInput, optionalArgs...)
}

func (c *Context) tmplTakeRoleIDWithError(target interface{}, roleID interface{}, optionalArgs ...interface{}) (string, error) {
	return c.tmplTakeRoleWithError(target, roleID, optionalArgs...)
}

func (c *Context) tmplTakeRoleNameWithError(target interface{}, roleName string, optionalArgs ...interface{}) (string, error) {
	return c.tmplTakeRoleWithError(target, roleName, optionalArgs...)
}

func (c *Context) removeRole(roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.IncreaseCheckStateLock() {
		return "", ErrTooManyCalls
	}

	var delay int
	if len(optionalArgs) > 0 {
		delay = ToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", errors.New("removeRole called on context with nil MemberState")
	}

	role := c.FindRole(roleInput)
	if role == nil {
		return "", fmt.Errorf("Role %v not found", roleInput)
	}

	if delay > 0 {
		err := scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
		if err != nil {
			return "", err
		}
	} else {
		err := common.RemoveRoleDS(c.MS, role.ID)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplRemoveRole(roleInput interface{}, optionalArgs ...interface{}) string {
	out, _ := c.removeRole(roleInput, optionalArgs...)
	return out
}

func (c *Context) tmplRemoveRoleID(roleID interface{}, optionalArgs ...interface{}) string {
	return c.tmplRemoveRole(roleID, optionalArgs...)
}

func (c *Context) tmplRemoveRoleName(roleName string, optionalArgs ...interface{}) string {
	return c.tmplRemoveRole(roleName, optionalArgs...)
}

func (c *Context) tmplRemoveRoleWithError(roleInput interface{}, optionalArgs ...interface{}) (string, error) {
	return c.removeRole(roleInput, optionalArgs...)
}

func (c *Context) tmplRemoveRoleIDWithError(roleID interface{}, optionalArgs ...interface{}) (string, error) {
	return c.tmplRemoveRoleWithError(roleID, optionalArgs...)
}

func (c *Context) tmplRemoveRoleNameWithError(roleName string, optionalArgs ...interface{}) (string, error) {
	return c.tmplRemoveRoleWithError(roleName, optionalArgs...)
}

func (c *Context) tmplSetRoles(target interface{}, roleSlice interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	targetMember, err := c.getMember(target)
	if err != nil {
		return "", err
	}

	targetID := targetMember.ID

	if c.IncreaseCheckCallCounter("set_roles"+discordgo.StrID(targetID), 1) {
		return "", errors.New("Too many calls to setRoles for specific user ID (max 1 / user)")
	}

	rSlice := reflect.ValueOf(roleSlice)
	switch rSlice.Kind() {
	case reflect.Slice, reflect.Array:
		// ok
	default:
		return "", errors.New("Value passed was not an array or slice")
	}

	if rSlice.Len() > 250 {
		return "", errors.New("Length of slice passed was > 250 (Discord role limit)")
	}

	roles := make([]int64, 0, rSlice.Len()+len(targetMember.Roles))
	for i := 0; i < rSlice.Len(); i++ {
		switch v := rSlice.Index(i).Interface().(type) {
		case string:
			role, err := c.tmplGetRole(v)
			if err == nil && role != nil {
				roles = append(roles, role.ID)
			}
		case int:
			roles = append(roles, int64(v))
		case int64:
			roles = append(roles, v)
		default:
			return "", errors.New("Could not convert slice to string slice")
		}
	}

	for _, r := range targetMember.Roles {
		role := c.GS.Role(true, r)
		if c.IncreaseCheckStateLock() {
			return "", ErrTooManyCalls
		}

		if role.Managed && !common.ContainsInt64Slice(roles, r) {
			roles = append(roles, r)
		}
	}

	err = common.BotSession.GuildMemberEdit(c.GS.ID, targetID, roles)
	if err != nil {
		return "", err
	}

	return "", nil
}

// Reactions functions

// Deletes reactions from a message either via reaction trigger or argument-set of emojis,
// needs channelID, messageID, userID, list of emojis - up to twenty
func (c *Context) tmplDelMessageReaction(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 4 {
			return reflect.Value{}, errors.New("Argumentos insuficientes (precisa de channelID, messageID, userID, emoji)")
		}

		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		if cID == 0 {
			return reflect.ValueOf("canal não existente"), nil
		}

		mID := ToInt64(args[1].Interface())
		uID := targetUserID(args[2].Interface())

		for _, reaction := range args[3:] {
			if c.IncreaseCheckGenericAPICall() {
				return reflect.Value{}, ErrTooManyAPICalls
			}

			if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionRemove(cID, mID, emojiArg(reaction.String()), uID); err != nil {
				return reflect.Value{}, err
			}
		}

		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)
}

func (c *Context) tmplDelAllMessageReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 2 {
			return reflect.Value{}, errors.New("Argumentos insuficientes (precisa de channelID, messageID, emojis[optional])")
		}

		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		if cID == 0 {
			return reflect.ValueOf("canal não existente"), nil
		}

		mID := ToInt64(args[1].Interface())

		if len(args) > 2 {
			for _, emoji := range args[2:] {
				if c.IncreaseCheckGenericAPICall() {
					return reflect.Value{}, ErrTooManyAPICalls
				}

				if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
					return reflect.Value{}, ErrTooManyCalls
				}

				if err := common.BotSession.MessageReactionRemoveEmoji(cID, mID, emojiArg(emoji.String())); err != nil {
					return reflect.Value{}, err
				}
			}

			return reflect.ValueOf(""), nil
		}

		if c.IncreaseCheckGenericAPICall() {
			return reflect.Value{}, ErrTooManyAPICalls
		}

		_ = common.BotSession.MessageReactionsRemoveAll(cID, mID)
		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)
}

func (c *Context) tmplAddReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if c.Msg == nil {
			return reflect.Value{}, nil
		}

		for _, reaction := range args {
			if c.IncreaseCheckCallCounter("add_reaction_trigger", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionAdd(c.Msg.ChannelID, c.Msg.ID, reaction.String()); err != nil {
				return reflect.Value{}, err
			}
		}

		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, true, values...)
}

func (c *Context) tmplAddResponseReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		for _, reaction := range args {
			if c.IncreaseCheckCallCounter("add_reaction_response", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			c.CurrentFrame.AddResponseReactionNames = append(c.CurrentFrame.AddResponseReactionNames, reaction.String())
		}

		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, true, values...)
}

func (c *Context) tmplAddMessageReactions(values ...reflect.Value) (reflect.Value, error) {
	f := func(args []reflect.Value) (reflect.Value, error) {
		if len(args) < 2 {
			return reflect.Value{}, errors.New("Not enough arguments (need channel and message-id)")
		}

		// cArg := args[0].Interface()
		var cArg interface{}
		if args[0].IsValid() {
			cArg = args[0].Interface()
		}

		cID := c.ChannelArg(cArg)
		mID := ToInt64(args[1].Interface())

		if cID == 0 {
			return reflect.ValueOf(""), nil
		}

		for i, reaction := range args {
			if i < 2 {
				continue
			}

			if c.IncreaseCheckCallCounter("add_reaction_message", 20) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionAdd(cID, mID, reaction.String()); err != nil {
				return reflect.Value{}, err
			}
		}

		return reflect.ValueOf(""), nil
	}

	return callVariadic(f, false, values...)
}

// Regex functions

func (c *Context) compileRegex(r string) (*regexp.Regexp, error) {
	if c.RegexCache == nil {
		c.RegexCache = make(map[string]*regexp.Regexp)
	}

	cached, ok := c.RegexCache[r]
	if ok {
		return cached, nil
	}

	if len(c.RegexCache) >= 10 {
		return nil, ErrTooManyCalls
	}

	compiled, err := regexp.Compile(r)
	if err != nil {
		return nil, err
	}

	c.RegexCache[r] = compiled

	return compiled, nil
}

func (c *Context) reFind(r string, s string) (string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return "", err
	}

	return compiled.FindString(s), nil
}

func (c *Context) reFindAll(r string, s string) ([]string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return nil, err
	}

	return compiled.FindAllString(s, 1000), nil
}

func (c *Context) reFindAllSubmatches(r string, s string) ([][]string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return nil, err
	}

	return compiled.FindAllStringSubmatch(s, 100), nil
}

func (c *Context) reReplace(r string, s string, repl string) (string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return "", err
	}

	return compiled.ReplaceAllString(s, repl), nil
}

func (c *Context) reSplit(r, s string, i int) ([]string, error) {
	compiled, err := c.compileRegex(r)
	if err != nil {
		return nil, err
	}

	return compiled.Split(s, i), nil
}

// Channels functions

func (c *Context) tmplEditChannelTopic(channel interface{}, newTopic string) (string, error) {
	if c.IncreaseCheckCallCounter("edit_channel", 10) {
		return "", ErrTooManyCalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return "", errors.New("Unknown channel")
	}

	if c.IncreaseCheckCallCounter("edit_channel_"+strconv.FormatInt(cID, 10), 2) {
		return "", ErrTooManyCalls
	}

	edit := &discordgo.ChannelEdit{
		Topic: newTopic,
	}

	_, err := common.BotSession.ChannelEditComplex(cID, edit)
	return "", err
}

func (c *Context) tmplEditChannelName(channel interface{}, newName string) (string, error) {
	if c.IncreaseCheckCallCounter("edit_channel", 10) {
		return "", ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return "", errors.New("Unknown channel")
	}

	if c.IncreaseCheckCallCounter("edit_channel_"+strconv.FormatInt(cID, 10), 2) {
		return "", ErrTooManyCalls
	}

	_, err := common.BotSession.ChannelEdit(cID, newName)
	return "", err
}

func (c *Context) tmplGetChannel(channel interface{}) (*CtxChannel, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	cID := c.ChannelArg(channel)
	if cID == 0 {
		return nil, nil //dont send an error , a nil output would indicate invalid/unknown channel
	}

	cstate := c.GS.ChannelCopy(true, cID)

	if cstate == nil {
		return nil, errors.New("Channel not in state")
	}

	return CtxChannelFromCS(cstate), nil

}

func (c *Context) tmplCreateChannel(name string, ctype interface{}) (*CtxChannel, error) {
	if c.IncreaseCheckCallCounterPremium("create_channel", 1, 2) {
		return nil, ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	if name == "" {
		return nil, errors.New("Channel name can not be empty")
	}

	ctypeInt := ToInt(ctype)
	switch ctypeInt {
	case 1, 3, 5, 6:
		return nil, errors.Errorf("Can not create a channel of type %d", ctypeInt)
	}

	dChannel, err := common.BotSession.GuildChannelCreate(c.GS.ID, name, discordgo.ChannelType(ctypeInt))
	if err != nil {
		return nil, err
	}

	return CtxChannelFromDGoChannel(dChannel), nil
}

func (c *Context) tmplCreateChannelComplex(channelDataArgs ...interface{}) (*CtxChannel, error) {
	if c.IncreaseCheckCallCounterPremium("create_channel", 1, 2) {
		return nil, ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	if len(channelDataArgs) < 1 {
		return nil, errors.New("Invalid channel data provided")
	}

	var m map[string]interface{}
	switch t := channelDataArgs[0].(type) {
	case SDict:
		m = t
	case map[string]interface{}:
		m = t
	default:
		dict, err := StringKeyDictionary(channelDataArgs...)
		if err != nil {
			return nil, err
		}
		m = dict
	}

	encoded, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var channelData *discordgo.GuildChannelCreateData
	err = json.Unmarshal(encoded, &channelData)
	if err != nil {
		return nil, err
	}

	if channelData == nil {
		return nil, errors.New("Nil channel data")
	}

	switch channelData.Type {
	case 1, 3, 5, 6:
		return nil, errors.Errorf("Can not create a channel of type %d", channelData.Type)
	}

	dChannel, err := common.BotSession.GuildChannelCreateComplex(c.GS.ID, *channelData)
	if err != nil {
		return nil, err
	}

	return CtxChannelFromDGoChannel(dChannel), nil
}

// TryCall functions

// tryCall attempts to call a non-builtin function by name.
// Modified version of https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/exec.go#L669
func (c *Context) tmplTryCall(name string, args ...reflect.Value) ([]interface{}, error) {
	// Probably will lead to oddities I haven't thought of, and no real reason to call tryCall using tryCall
	if name == "tryCall" {
		return nil, errors.New("can't call tryCall function using tryCall")
	}

	fun, ok := c.ContextFuncs[name]
	if !ok {
		// try looking up the name from StandardFuncMap
		fun, ok = StandardFuncMap[name]
		if !ok {
			return nil, errors.Errorf("function %q not found", name)
		}
	}

	funVal := reflect.ValueOf(fun)
	typ := funVal.Type()

	numIn := len(args)
	numFixed := len(args)
	if typ.IsVariadic() {
		numFixed = typ.NumIn() - 1 // last arg is the variadic one.
		if numIn < numFixed {
			return nil, errors.Errorf("wrong number of args for %s: want at least %d got %d", name, typ.NumIn()-1, len(args))
		}
	} else if numIn != typ.NumIn() {
		return nil, errors.Errorf("wrong number of args for %s: want %d got %d", name, typ.NumIn(), numIn)
	}

	if !template.GoodFunc(typ) {
		return nil, errors.Errorf("can't call function %q with %d results", name, typ.NumOut())
	}

	argv := make([]reflect.Value, numIn)
	var i int
	// Validate fixed args.
	for ; i < numFixed && i < len(args); i++ {
		val, err := validateType(args[i], typ.In(i))
		if err != nil {
			return nil, err
		}

		argv[i] = val
	}

	// Now the ... args.
	if typ.IsVariadic() {
		argType := typ.In(typ.NumIn() - 1).Elem() // Argument is a slice.
		for ; i < len(args); i++ {
			val, err := validateType(args[i], argType)
			if err != nil {
				return nil, err
			}

			argv[i] = val
		}
	}

	v, err, _ := template.SafeCall(funVal, argv)
	if err != nil {
		return []interface{}{nil, err}, nil
	}

	if v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}

	return []interface{}{v, nil}, nil
}

// validateType guarantees that the value is valid and assignable to the type.
// Taken from https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/exec.go#L746
func validateType(value reflect.Value, typ reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if typ == nil {
			// An untyped nil interface{}. Accept as a proper nil value.
			return reflect.ValueOf(nil), nil
		}

		if template.CanBeNil(typ) {
			// Like above, but use the zero value of the non-nil type.
			return reflect.Zero(typ), nil
		}

		return reflect.Value{}, errors.Errorf("invalid value; expected %s", typ)
	}

	if typ == reflectValueType && value.Type() != typ {
		return reflect.ValueOf(value), nil
	}

	if typ != nil && !value.Type().AssignableTo(typ) {
		if value.Kind() == reflect.Interface && !value.IsNil() {
			value = value.Elem()
			if value.Type().AssignableTo(typ) {
				return value, nil
			}
			// fallthrough
		}

		// Does one dereference or indirection work? We could do more, as we
		// do with method receivers, but that gets messy and method receivers
		// are much more constrained, so it makes more sense there than here.
		// Besides, one is almost always all you need.
		switch {
		case value.Kind() == reflect.Ptr && value.Type().Elem().AssignableTo(typ):
			value = value.Elem()
			if !value.IsValid() {
				return reflect.Value{}, errors.Errorf("dereference of nil pointer of type %s", typ)
			}
		case reflect.PtrTo(value.Type()).AssignableTo(typ) && value.CanAddr():
			value = value.Addr()
		default:
			return reflect.Value{}, errors.Errorf("wrong type for value; expected %s; got %s", typ, value.Type())
		}
	}

	return value, nil
}

// Standardize functions

type StdDepth struct {
	depth int
}

func newStdDepth() *StdDepth {
	return &StdDepth{
		depth: 0,
	}
}

func (c *Context) tmplStandardize(input interface{}) interface{} {
	depth := newStdDepth()
	return depth.StdInit(input)
}

func (sd *StdDepth) StdInit(input interface{}) interface{} {
	val := reflect.ValueOf(input)
	switch val.Kind() {
	case reflect.Map:
		switch t := input.(type) {
		case Dict:
			return sd.StdMap(t)
		case SDict:
			return sd.StdStringMap(t)
		case map[string]interface{}:
			return sd.StdStringMap(t)
		default:
			return sd.StdMap(t)
		}
	default:
		switch t := input.(type) {
		case []interface{}:
			return sd.StdSlice(t)
		case *time.Time:
			return *t
		default:
			return t
		}
	}
}

func (sd *StdDepth) StdMap(input interface{}) interface{} {
	if sd.depth >= 1000 {
		return input
	}
	sd.depth++
	out := make(Dict)
	val := reflect.ValueOf(input)
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			switch t := v.Interface().(type) {
			case map[interface{}]interface{}:
				_, err := out.Set(k.Interface(), sd.StdMap(t))
				if err != nil {
					return input
				}
			case map[string]interface{}:
				_, err := out.Set(k.Interface(), sd.StdStringMap(t))
				if err != nil {
					return input
				}
			case []interface{}:
				_, err := out.Set(k.Interface(), sd.StdSlice(t))
				if err != nil {
					return input
				}
			default:
				_, err := out.Set(k.Interface(), sd.StdInit(t))
				if err != nil {
					return input
				}
			}
		}
		return out
	}
	return nil
}

func (sd *StdDepth) StdStringMap(input interface{}) interface{} {
	if sd.depth >= 1000 {
		return input
	}
	sd.depth++
	out := make(SDict)
	val := reflect.ValueOf(input)
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			v := val.MapIndex(k)
			switch t := v.Interface().(type) {
			case map[interface{}]interface{}:
				_, err := out.Set(k.Interface().(string), sd.StdMap(t))
				if err != nil {
					return input
				}
			case map[string]interface{}:
				_, err := out.Set(k.Interface().(string), sd.StdStringMap(t))
				if err != nil {
					return input
				}
			case []interface{}:
				_, err := out.Set(k.Interface().(string), sd.StdSlice(t))
				if err != nil {
					return input
				}
			default:
				_, err := out.Set(k.Interface().(string), sd.StdInit(t))
				if err != nil {
					return input
				}
			}
		}
		return out
	}
	return nil
}

func (sd *StdDepth) StdSlice(input []interface{}) interface{} {
	if sd.depth >= 1000 {
		return input
	}
	sd.depth++
	var out Slice
	for _, v := range input {
		out = append(out, sd.StdInit(v))
	}

	return out
}

// Misc functions

func (c *Context) tmplGetMember(target interface{}) (*CtxMember, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	member, err := c.getMember(target)
	if err != nil {
		return nil, err
	}

	return CtxMemberFromMS(member), nil
}

func (c *Context) tmplCurrentUserCreated() time.Time {
	t := bot.SnowflakeToTime(c.MS.ID)
	return t
}

func (c *Context) tmplCurrentUserAgeHuman() string {
	t := bot.SnowflakeToTime(c.MS.ID)

	humanized := common.HumanizeDuration(common.DurationPrecisionHours, time.Since(t))
	if humanized == "" {
		humanized = "Menos de uma hora!"
	}

	return humanized
}

func (c *Context) tmplCurrentUserAgeMinutes() int {
	t := bot.SnowflakeToTime(c.MS.ID)
	d := time.Since(t)

	return int(d.Seconds() / 60)
}

func (c *Context) tmplUserCreated(id interface{}) *time.Time {
	if c.IncreaseCheckGenericAPICall() {
		return nil
	}

	ms, _ := c.getMember(id)
	if ms == nil {
		return nil
	}

	t := bot.SnowflakeToTime(ms.ID)
	return &t
}

func (c *Context) tmplUserAgeHuman(id interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	ms, _ := c.getMember(id)
	if ms == nil {
		return ""
	}

	t := bot.SnowflakeToTime(ms.ID)

	humanized := common.HumanizeDuration(common.DurationPrecisionHours, time.Since(t))
	if humanized == "" {
		humanized = "Menos de uma hora!"
	}

	return humanized
}

func (c *Context) tmplUserAgeMinutes(id interface{}) int {
	if c.IncreaseCheckGenericAPICall() {
		return 0
	}

	ms, _ := c.getMember(id)
	if ms == nil {
		return 0
	}

	t := bot.SnowflakeToTime(ms.ID)
	d := time.Since(t)

	return int(d.Seconds() / 60)
}

func (c *Context) tmplSleep(duration interface{}) (string, error) {
	seconds := ToInt(duration)
	if c.secondsSlept+seconds > 60 || seconds < 1 {
		return "", errors.New("can sleep for max 60 seconds combined")
	}

	c.secondsSlept += seconds
	time.Sleep(time.Duration(seconds) * time.Second)

	return "", nil
}

func (c *Context) tmplOnlineCount() (int, error) {
	if c.IncreaseCheckCallCounter("online_users", 1) {
		return 0, ErrTooManyCalls
	}

	if c.IncreaseCheckStateLock() {
		return 0, ErrTooManyCalls
	}

	var online int
	c.GS.RLock()
	for _, v := range c.GS.Members {
		if v.PresenceSet && v.PresenceStatus != dstate.StatusOffline {
			online++
		}
	}
	c.GS.RUnlock()

	return online, nil
}

func (c *Context) tmplOnlineCountBots() (int, error) {
	if c.IncreaseCheckCallCounter("online_bots", 1) {
		return 0, ErrTooManyCalls
	}

	if c.IncreaseCheckStateLock() {
		return 0, ErrTooManyCalls
	}

	var botCount int
	c.GS.RLock()
	for _, v := range c.GS.Members {
		if v.Bot && v.PresenceSet && v.PresenceStatus != dstate.StatusOffline {
			botCount++
		}
	}
	c.GS.RUnlock()

	return botCount, nil
}

func (c *Context) tmplEditNickname(nick string) (string, error) {
	if c.IncreaseCheckCallCounter("edit_nick", 2) {
		return "", ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	if c.MS == nil {
		return "", nil
	}

	if strings.Compare(c.MS.Nick, nick) == 0 {
		return "", nil
	}

	err := common.BotSession.GuildMemberNickname(c.GS.ID, c.MS.ID, nick)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Context) tmplEditTargetNickName(target int64, nick string) (string, error) {
	if !bot.IsGuildWhiteListed(c.GS.ID) {
		return "", errors.New("Esse server não tem permissão para usar a função editTargetNickaname")
	}

	if c.IncreaseCheckCallCounter("edit_nick", 2) {
		return "", ErrTooManyCalls
	}

	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	ms, err := c.getMember(target)
	if err != nil {
		return "", err
	}

	if strings.Compare(ms.Nick, nick) == 0 {
		return "", nil
	}

	err = common.BotSession.GuildMemberNickname(c.GS.ID, ms.ID, nick)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Context) tmplSort(slice []interface{}, sortargs ...interface{}) (interface{}, error) {
	if c.IncreaseCheckCallCounterPremium("sortfuncs", 1, 3) {
		return "", ErrTooManyCalls
	}

	var dict SDict
	var err error
	switch len(sortargs) {
	case 0:
		input := make(map[string]interface{}, 3)
		input["reverse"] = false
		input["subslices"] = false
		input["emptyslices"] = false

		dict, err = StringKeyDictionary(input)
		if err != nil {
			return "", err
		}
	case 1:
		dict, err = StringKeyDictionary(sortargs[0])
		if err != nil {
			return "", err
		}
	default:
		dict, err = StringKeyDictionary(sortargs...)
		if err != nil {
			return "", err
		}
	}

	var numberSlice, stringSlice, timeSlice, csliceSlice, mapSlice, defaultSlice, outputSlice []interface{}

	for _, v := range slice {
		switch t := v.(type) {
		case int, int64, float64:
			numberSlice = append(numberSlice, t)
		case string:
			stringSlice = append(stringSlice, t)
		case time.Time:
			timeSlice = append(timeSlice, t)
		case *time.Time:
			if t != nil {
				timeSlice = append(timeSlice, *t)
			}
		default:
			v := reflect.ValueOf(t)
			switch v.Kind() {
			case reflect.Slice:
				csliceSlice = append(csliceSlice, t)
			case reflect.Map:
				mapSlice = append(mapSlice, t)
			default:
				defaultSlice = append(defaultSlice, t)
			}
		}
	}

	if dict.Get("reverse") == true {
		sort.Slice(numberSlice, func(i, j int) bool { return ToFloat64(numberSlice[i]) > ToFloat64(numberSlice[j]) })
		sort.Slice(stringSlice, func(i, j int) bool { return ToString(stringSlice[i]) > ToString(stringSlice[j]) })
		sort.Slice(timeSlice, func(i, j int) bool { return timeSlice[i].(time.Time).Before(timeSlice[j].(time.Time)) })
		sort.Slice(csliceSlice, func(i, j int) bool { return getLen(csliceSlice[i]) > getLen(csliceSlice[j]) })
		sort.Slice(mapSlice, func(i, j int) bool { return getLen(mapSlice[i]) > getLen(mapSlice[j]) })
	} else {
		sort.Slice(numberSlice, func(i, j int) bool { return ToFloat64(numberSlice[i]) < ToFloat64(numberSlice[j]) })
		sort.Slice(stringSlice, func(i, j int) bool { return ToString(stringSlice[i]) < ToString(stringSlice[j]) })
		sort.Slice(timeSlice, func(i, j int) bool { return timeSlice[j].(time.Time).Before(timeSlice[i].(time.Time)) })
		sort.Slice(csliceSlice, func(i, j int) bool { return getLen(csliceSlice[i]) < getLen(csliceSlice[j]) })
		sort.Slice(mapSlice, func(i, j int) bool { return getLen(mapSlice[i]) < getLen(mapSlice[j]) })
	}

	if dict.Get("subslices") == true {
		if dict.Get("emptyslices") == true {
			outputSlice = append(outputSlice, numberSlice, stringSlice, timeSlice, csliceSlice, mapSlice, defaultSlice)
		} else {
			if len(numberSlice) > 0 {
				outputSlice = append(outputSlice, numberSlice)
			}

			if len(stringSlice) > 0 {
				outputSlice = append(outputSlice, stringSlice)
			}

			if len(timeSlice) > 0 {
				outputSlice = append(outputSlice, timeSlice)
			}

			if len(csliceSlice) > 0 {
				outputSlice = append(outputSlice, csliceSlice)
			}

			if len(mapSlice) > 0 {
				outputSlice = append(outputSlice, mapSlice)
			}

			if len(defaultSlice) > 0 {
				outputSlice = append(outputSlice, defaultSlice)
			}
		}
	} else {
		outputSlice = append(outputSlice, numberSlice...)
		outputSlice = append(outputSlice, stringSlice...)
		outputSlice = append(outputSlice, timeSlice...)
		outputSlice = append(outputSlice, csliceSlice...)
		outputSlice = append(outputSlice, mapSlice...)
		outputSlice = append(outputSlice, defaultSlice...)
	}

	return outputSlice, nil
}

func (c *Context) tmplGeneratePerms(perms ...interface{}) (string, error) {
	l := len(perms)
	var permSlice []interface{}
	var err error

	switch {
	case l < 1:
		return "", errors.New("Not enough arguments provided")
	case l == 1:
		switch t := perms[0].(type) {
		case int, int64, float64:
			return ToString(t), nil
		case []interface{}:
			permSlice = t
		case Slice:
			permSlice = t
		case string:
			var strOut int64
			if permMap[t] != 0 {
				strOut |= permMap[t]
				return ToString(strOut), nil
			}

			return "", errors.New("Invalid perm provided.")
		}
	default:
		permSlice, err = CreateSlice(perms...)
		if err != nil {
			return "", err
		}
	}

	var perm int64
	for _, v := range permSlice {
		vStr, ok := v.(string)
		if !ok {
			return "", errors.New("Non string value found on slice")
		}

		if permMap[vStr] != 0 {
			perm |= permMap[vStr]
		}
	}

	return ToString(perm), nil
}
