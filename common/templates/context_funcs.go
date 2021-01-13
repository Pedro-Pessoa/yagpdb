package templates

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/scheduledevents2"
)

var (
	ErrTooManyCalls    = errors.New("Too many calls to this function")
	ErrTooManyAPICalls = errors.New("Too many potential discord api calls function")

	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

func (c *Context) buildDM(gName string, s ...interface{}) *discordgo.MessageSend {
	info := fmt.Sprintf("DM enviada pelo servidor **%s**", gName)
	msgSend := &discordgo.MessageSend{
		AllowedMentions: discordgo.AllowedMentions{
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
	if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 1) || c.MS == nil {
		return ""
	}

	c.GS.RLock()
	memberID, gName := c.MS.ID, c.GS.Guild.Name
	c.GS.RUnlock()

	msgSend := c.buildDM(gName, s...)
	if msgSend == nil {
		return ""
	}

	channel, err := common.BotSession.UserChannelCreate(memberID)
	if err != nil {
		return ""
	}
	_, _ = common.BotSession.ChannelMessageSendComplex(channel.ID, msgSend)
	return ""
}

func (c *Context) tmplSendTargetDM(target interface{}, s ...interface{}) string {
	if bot.IsGuildWhiteListed(c.GS.ID) {
		if len(s) < 1 || c.IncreaseCheckCallCounter("send_dm", 1) {
			return ""
		}

		targetID := targetUserID(target)
		if targetID == 0 {
			return ""
		}

		ts, err := bot.GetMember(c.GS.ID, targetID)
		if err != nil {
			return ""
		}

		msgSend := c.buildDM("", s...)
		if msgSend == nil {
			return ""
		}

		channel, err := common.BotSession.UserChannelCreate(ts.ID)
		if err != nil {
			return ""
		}
		_, _ = common.BotSession.ChannelMessageSendComplex(channel.ID, msgSend)
	}

	return ""
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

func (c *Context) tmplSendTemplateDM(name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(nil, true, false, name, data...)
}

func (c *Context) tmplSendTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(channel, false, false, name, data...)
}

func (c *Context) tmplExecTemplate(channel interface{}, name string, data ...interface{}) (interface{}, error) {
	return c.sendNestedTemplate(channel, false, true, name, data...)
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

func (c *Context) tmplSendMessage(filterSpecialMentions bool, returnID bool) func(channel interface{}, msg interface{}) interface{} {
	parseMentions := []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers}
	if !filterSpecialMentions {
		parseMentions = append(parseMentions, discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeEveryone)
	}

	return func(channel interface{}, msg interface{}) interface{} {
		if c.IncreaseCheckGenericAPICall() {
			return ""
		}

		cid := c.ChannelArg(channel)
		if cid == 0 {
			return ""
		}

		isDM := cid != c.ChannelArgNoDM(channel)
		c.GS.RLock()
		info := fmt.Sprintf("DM enviada pelo servidor **%s**", c.GS.Guild.Name)
		c.GS.RUnlock()
		WL := bot.IsGuildWhiteListed(c.GS.ID)

		var m *discordgo.Message
		msgSend := &discordgo.MessageSend{
			AllowedMentions: discordgo.AllowedMentions{
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
				msgSend.AllowedMentions = discordgo.AllowedMentions{Parse: parseMentions}
			}

			if isDM && !WL {
				if typedMsg.Embed != nil {
					typedMsg.Embed.Footer.Text = info
				} else {
					typedMsg.Content = info + "\n" + typedMsg.Content
				}
			}
		default:
			if isDM && !WL {
				msgSend.Content = info + "\n" + fmt.Sprint(msg)
			} else {
				msgSend.Content = fmt.Sprint(msg)
			}
		}

		m, err = common.BotSession.ChannelMessageSendComplex(cid, msgSend)

		if err == nil && returnID {
			return m.ID
		}

		return ""
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
			temp := fmt.Sprint(msg)
			msgEdit.Content = &temp
		}

		if !filterSpecialMentions {
			msgEdit.AllowedMentions = &discordgo.AllowedMentions{
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

func (c *Context) tmplMentionEveryone() string {
	c.CurrentFrame.MentionEveryone = true
	return "@everyone"
}

func (c *Context) tmplMentionHere() string {
	c.CurrentFrame.MentionHere = true
	return "@here"
}

func (c *Context) tmplMentionRoleID(roleID interface{}) string {
	if c.IncreaseCheckStateLock() {
		return ""
	}

	role := ToInt64(roleID)
	if role == int64(0) {
		return ""
	}

	r := c.GS.RoleCopy(true, role)
	if r == nil {
		return "(role not found)"
	}

	if common.ContainsInt64Slice(c.CurrentFrame.MentionRoles, role) {
		return "<@&" + discordgo.StrID(role) + ">"
	}

	c.CurrentFrame.MentionRoles = append(c.CurrentFrame.MentionRoles, role)
	return "<@&" + discordgo.StrID(role) + ">"
}

func (c *Context) tmplMentionRoleName(role string) string {
	if c.IncreaseCheckStateLock() {
		return ""
	}

	var found *discordgo.Role
	c.GS.RLock()
	for _, r := range c.GS.Guild.Roles {
		if strings.EqualFold(r.Name, role) {
			if !common.ContainsInt64Slice(c.CurrentFrame.MentionRoles, r.ID) {
				c.CurrentFrame.MentionRoles = append(c.CurrentFrame.MentionRoles, r.ID)
				found = r
			}
			break
		}
	}
	c.GS.RUnlock()
	if found == nil {
		return "(role not found)"
	}

	return "<@&" + discordgo.StrID(found.ID) + ">"
}

func (c *Context) tmplHasRoleID(roleID interface{}) bool {
	role := ToInt64(roleID)
	if role == 0 {
		return false
	}

	return common.ContainsInt64Slice(c.MS.Roles, role)
}

func (c *Context) tmplHasRoleName(name string) (bool, error) {
	if c.IncreaseCheckStateLock() {
		return false, ErrTooManyCalls
	}

	role := c.findRoleByName(name)
	if role == nil {
		return false, nil
	}

	if common.ContainsInt64Slice(c.MS.Roles, role.ID) {
		return true, nil
	}

	// Role not found, default to false
	return false, nil
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

func (c *Context) tmplTargetHasRoleID(target interface{}, roleID interface{}) bool {
	if c.IncreaseCheckStateLock() {
		return false
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return false
	}

	ts, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return false
	}

	role := ToInt64(roleID)
	if role == 0 {
		return false
	}

	return common.ContainsInt64Slice(ts.Roles, role)

}

func (c *Context) tmplTargetHasRoleName(target interface{}, name string) bool {
	if c.IncreaseCheckStateLock() {
		return false
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return false
	}

	ts, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return false
	}

	role := c.findRoleByName(name)
	if role == nil {
		return false
	}

	return common.ContainsInt64Slice(ts.Roles, role.ID)

}

func (c *Context) tmplGiveRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := ToInt64(roleID)
	if role == 0 {
		return ""
	}

	// Check to see if we can save a API request here, if this isn't delayed
	if delay <= 0 {
		c.GS.RLock()
		ms := c.GS.Member(false, targetID)
		hasRole := true
		if ms != nil && ms.MemberSet {
			hasRole = common.ContainsInt64Slice(ms.Roles, role)
		}
		c.GS.RUnlock()

		if !hasRole {
			return ""
		}
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, targetID, role, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		_ = common.BotSession.GuildMemberRoleAdd(c.GS.ID, targetID, role)
	}

	return ""
}

func (c *Context) tmplGiveRoleName(target interface{}, name string, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "no role by the name of " + name + " found"
	}

	// Maybe save an api request
	if delay <= 0 {
		ms := c.GS.Member(false, targetID)
		if ms != nil {
			if common.ContainsInt64Slice(ms.Roles, role.ID) {
				return ""
			}
		}
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		_ = common.BotSession.GuildMemberRoleAdd(c.GS.ID, targetID, role.ID)
	}

	return ""
}

func (c *Context) tmplTakeRoleID(target interface{}, roleID interface{}, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := ToInt64(roleID)
	if role == 0 {
		return ""
	}

	// Check to see if we can save a API request here, if this isn't delayed
	if delay <= 0 {
		c.GS.RLock()
		ms := c.GS.Member(false, targetID)
		hasRole := true
		if ms != nil && ms.MemberSet {
			hasRole = common.ContainsInt64Slice(ms.Roles, role)
		}
		c.GS.RUnlock()

		if !hasRole {
			return ""
		}
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, targetID, role, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		_ = common.BotSession.GuildMemberRoleRemove(c.GS.ID, targetID, role)
	}

	return ""
}

func (c *Context) tmplTakeRoleName(target interface{}, name string, optionalArgs ...interface{}) string {
	if c.IncreaseCheckGenericAPICall() {
		return ""
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return ""
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "no role by the name of " + name + " found"
	}

	// Maybe save an api request
	if delay <= 0 {
		ms := c.GS.Member(false, targetID)
		if ms != nil {
			if common.ContainsInt64Slice(ms.Roles, role.ID) {
				return ""
			}
		}
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, targetID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		_ = common.BotSession.GuildMemberRoleRemove(c.GS.ID, targetID, role.ID)
	}

	return ""
}

func (c *Context) tmplAddRoleID(role interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	rid := ToInt64(role)
	if rid == 0 {
		return "", errors.New("No role id specified")
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, c.MS.ID, rid, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.AddRoleDS(c.MS, rid); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplAddRoleName(name string, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "", errors.New("No Role with name " + name + " found")
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleAddRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.AddRoleDS(c.MS, role.ID); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplRemoveRoleID(role interface{}, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	rid := ToInt64(role)
	if rid == 0 {
		return "", errors.New("No role id specified")
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, c.MS.ID, rid, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.RemoveRoleDS(c.MS, rid); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) tmplRemoveRoleName(name string, optionalArgs ...interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	delay := 0
	if len(optionalArgs) > 0 {
		delay = tmplToInt(optionalArgs[0])
	}

	if c.MS == nil {
		return "", nil
	}

	role := c.findRoleByName(name)
	if role == nil {
		return "", errors.New("No Role with name " + name + " found")
	}

	if delay > 0 {
		_ = scheduledevents2.ScheduleRemoveRole(context.Background(), c.GS.ID, c.MS.ID, role.ID, time.Now().Add(time.Second*time.Duration(delay)))
	} else {
		if err := common.RemoveRoleDS(c.MS, role.ID); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (c *Context) findRoleByID(id int64) *discordgo.Role {
	c.GS.RLock()
	defer c.GS.RUnlock()

	for _, r := range c.GS.Guild.Roles {
		if r.ID == id {
			return r
		}
	}

	return nil
}

func (c *Context) findRoleByName(name string) *discordgo.Role {
	c.GS.RLock()
	defer c.GS.RUnlock()

	for _, r := range c.GS.Guild.Roles {
		if strings.EqualFold(r.Name, name) {
			return r
		}
	}

	return nil
}

func (c *Context) tmplDelResponse(args ...interface{}) string {
	dur := 10
	if len(args) > 0 {
		dur = tmplToInt(args[0])
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
		dur = tmplToInt(args[0])
	}

	if dur > 86400 {
		dur = 86400
	}

	MaybeScheduledDeleteMessage(c.GS.ID, cID, mID, dur)

	return ""
}

//Deletes reactions from a message either via reaction trigger or argument-set of emojis,
//needs channelID, messageID, userID, list of emojis - up to twenty
//can be run once per CC.
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
			if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
				return reflect.Value{}, ErrTooManyCalls
			}

			if err := common.BotSession.MessageReactionRemove(cID, mID, reaction.String(), uID); err != nil {
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
				if c.IncreaseCheckCallCounter("del_reaction_message", 10) {
					return reflect.Value{}, ErrTooManyCalls
				}

				if err := common.BotSession.MessageReactionRemoveEmoji(cID, mID, emoji.String()); err != nil {
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

func (c *Context) tmplGetMessage(channel, msgID interface{}) (*discordgo.Message, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	cID := c.ChannelArgNoDM(channel)
	if cID == 0 {
		return nil, nil
	}

	mID := ToInt64(msgID)

	message, _ := common.BotSession.ChannelMessage(cID, mID)
	return message, nil
}

func (c *Context) tmplGetMember(target interface{}) (*CtxMember, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	mID := targetUserID(target)
	if mID == 0 {
		return nil, nil
	}

	member, _ := bot.GetMember(c.GS.ID, mID)
	if member == nil {
		return nil, nil
	}

	return CtxMemberFromMS(member), nil
}

func (c *Context) tmplGetRole(r interface{}) (*discordgo.Role, error) {
	if c.IncreaseCheckGenericAPICall() {
		return nil, ErrTooManyAPICalls
	}

	switch t := r.(type) {
	case int, int64:
		return c.findRoleByID(ToInt64(t)), nil
	case string:
		parsed, err := strconv.ParseInt(t, 10, 64)
		if err == nil {
			return c.findRoleByID(parsed), nil
		}

		if strings.HasPrefix(t, "<@&") && strings.HasSuffix(t, ">") && (len(t) > 4) {
			parsedMention, err := strconv.ParseInt(t[2:len(t)-1], 10, 64)
			if err == nil {
				return c.findRoleByID(parsedMention), nil
			}
		}

		return c.findRoleByName(t), nil
	default:
		return nil, nil
	}
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

func (c *Context) tmplCurrentUserCreated() time.Time {
	t := bot.SnowflakeToTime(c.MS.ID)
	return t
}

func (c *Context) tmplUserAgeHuman(id int64) *string {
	ms := c.getMember(id)
	if ms == nil {
		return nil
	}

	t := bot.SnowflakeToTime(ms.ID)

	humanized := common.HumanizeDuration(common.DurationPrecisionHours, time.Since(t))
	if humanized == "" {
		humanized = "Menos de uma hora!"
	}

	return &humanized
}

func (c *Context) tmplUserAgeMinutes(id int64) *int {
	ms := c.getMember(id)
	if ms == nil {
		return nil
	}

	t := bot.SnowflakeToTime(ms.ID)
	d := time.Since(t)

	out := int(d.Seconds() / 60)
	return &out
}

func (c *Context) tmplUserCreated(id int64) *time.Time {
	ms := c.getMember(id)
	if ms == nil {
		return nil
	}

	t := bot.SnowflakeToTime(ms.ID)
	return &t
}

func (c *Context) getMember(id int64) *dstate.MemberState {
	targetID := targetUserID(id)
	if targetID == 0 {
		return nil
	}

	ms, err := bot.GetMember(c.GS.ID, targetID)
	if err != nil {
		return nil
	}

	return ms
}

func (c *Context) tmplSleep(duration interface{}) (string, error) {
	seconds := tmplToInt(duration)
	if c.secondsSlept+seconds > 60 || seconds < 1 {
		return "", errors.New("can sleep for max 60 seconds combined")
	}

	c.secondsSlept += seconds
	time.Sleep(time.Duration(seconds) * time.Second)
	return "", nil
}

func (c *Context) compileRegex(r string) (*regexp.Regexp, error) {
	if c.RegexCache == nil {
		c.RegexCache = make(map[string]*regexp.Regexp)
	}

	cached, ok := c.RegexCache[r]
	if ok {
		return cached, nil
	}

	if len(c.RegexCache) >= 10 {
		return nil, ErrTooManyAPICalls
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

func (c *Context) tmplEditChannelName(channel interface{}, newName string) (string, error) {
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

	_, err := common.BotSession.ChannelEdit(cID, newName)
	return "", err
}

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

func (c *Context) tmplOnlineCount() (int, error) {
	if c.IncreaseCheckCallCounter("online_users", 1) {
		return 0, ErrTooManyCalls
	}

	online := 0
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

	botCount := 0
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

	ms := c.getMember(target)
	if ms == nil {
		return "", errors.New("Targeted user not found to change their nickname.")
	}

	if strings.Compare(ms.Nick, nick) == 0 {
		return "", nil
	}

	err := common.BotSession.GuildMemberNickname(c.GS.ID, ms.ID, nick)
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
	switch len(sortargs) {
	case 0:
		input := make(map[string]interface{}, 3)
		input["Reverse"] = false
		input["Subslices"] = false
		input["Emptyslices"] = false
		sdict, err := StringKeyDictionary(input)
		if err != nil {
			return "", err
		}
		dict = sdict
	case 1:
		sdict, err := StringKeyDictionary(sortargs[0])
		if err != nil {
			return "", err
		}
		dict = sdict
	default:
		return "", errors.New("Too many args.")
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
			timeSlice = append(timeSlice, *t)
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

	if dict.Get("Reverse") == true {
		sort.Slice(numberSlice, func(i, j int) bool { return ToFloat64(numberSlice[i]) > ToFloat64(numberSlice[j]) })
		sort.Slice(stringSlice, func(i, j int) bool { return getString(stringSlice[i]) > getString(stringSlice[j]) })
		sort.Slice(timeSlice, func(i, j int) bool { return timeSlice[i].(time.Time).Before(timeSlice[j].(time.Time)) })
		sort.Slice(csliceSlice, func(i, j int) bool { return getLen(csliceSlice[i]) > getLen(csliceSlice[j]) })
		sort.Slice(mapSlice, func(i, j int) bool { return getLen(mapSlice[i]) > getLen(mapSlice[j]) })
	} else {
		sort.Slice(numberSlice, func(i, j int) bool { return ToFloat64(numberSlice[i]) < ToFloat64(numberSlice[j]) })
		sort.Slice(stringSlice, func(i, j int) bool { return getString(stringSlice[i]) < getString(stringSlice[j]) })
		sort.Slice(timeSlice, func(i, j int) bool { return timeSlice[j].(time.Time).Before(timeSlice[i].(time.Time)) })
		sort.Slice(csliceSlice, func(i, j int) bool { return getLen(csliceSlice[i]) < getLen(csliceSlice[j]) })
		sort.Slice(mapSlice, func(i, j int) bool { return getLen(mapSlice[i]) < getLen(mapSlice[j]) })
	}

	if dict.Get("Subslices") == true {
		if dict.Get("Emptyslices") == true {
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

func getLen(from interface{}) int {
	v := reflect.ValueOf(from)
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len()
	default:
		return 0
	}
}

func getString(from interface{}) string {
	v := reflect.ValueOf(from)
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintln(from)
	default:
		return ""
	}
}

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
			return nil, fmt.Errorf("function %q not found", name)
		}
	}

	funVal := reflect.ValueOf(fun)
	typ := funVal.Type()

	numIn := len(args)
	numFixed := len(args)
	if typ.IsVariadic() {
		numFixed = typ.NumIn() - 1 // last arg is the variadic one.
		if numIn < numFixed {
			return nil, fmt.Errorf("wrong number of args for %s: want at least %d got %d", name, typ.NumIn()-1, len(args))
		}
	} else if numIn != typ.NumIn() {
		return nil, fmt.Errorf("wrong number of args for %s: want %d got %d", name, typ.NumIn(), numIn)
	}

	if !goodFunc(typ) {
		return nil, fmt.Errorf("can't call function %q with %d results", name, typ.NumOut())
	}

	argv := make([]reflect.Value, numIn)
	i := 0
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

	v, err := safeCall(funVal, argv)
	if err != nil {
		return []interface{}{nil, err}, nil
	}

	if v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}

	return []interface{}{v, nil}, nil
}

// safeCall runs fun.Call(args), and returns the resulting value and error, if
// any. If the call panics, the panic value is returned as an error.
// Taken from https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/funcs.go#L355
func safeCall(fun reflect.Value, args []reflect.Value) (val reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	ret := fun.Call(args)
	if len(ret) == 2 && !ret[1].IsNil() {
		return ret[0], ret[1].Interface().(error)
	}
	return ret[0], nil
}

// goodFunc reports whether the function or method has the right result signature.
// Taken from https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/funcs.go#L110
func goodFunc(typ reflect.Type) bool {
	// We allow functions with 1 result or 2 results where the second is an error.
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
// Taken from https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/exec.go#L735
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	case reflect.Struct:
		return typ == reflectValueType
	}
	return false
}

// validateType guarantees that the value is valid and assignable to the type.
// Taken from https://github.com/golang/go/blob/3b2a578166bdedd94110698c971ba8990771eb89/src/text/template/exec.go#L746
func validateType(value reflect.Value, typ reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if typ == nil {
			// An untyped nil interface{}. Accept as a proper nil value.
			return reflect.ValueOf(nil), nil
		}
		if canBeNil(typ) {
			// Like above, but use the zero value of the non-nil type.
			return reflect.Zero(typ), nil
		}
		return reflect.Value{}, fmt.Errorf("invalid value; expected %s", typ)
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
				return reflect.Value{}, fmt.Errorf("dereference of nil pointer of type %s", typ)
			}
		case reflect.PtrTo(value.Type()).AssignableTo(typ) && value.CanAddr():
			value = value.Addr()
		default:
			return reflect.Value{}, fmt.Errorf("wrong type for value; expected %s; got %s", typ, value.Type())
		}
	}
	return value, nil
}

type StdDepth struct {
	depth int
}

func newStdDepth() *StdDepth {
	return &StdDepth{
		depth: 0,
	}
}

func (sd *StdDepth) Add() {
	sd.depth += 1
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
		case Dict, SDict:
			return t
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
	sd.Add()
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
	sd.Add()
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
	sd.Add()
	var out Slice
	for _, v := range input {
		out = append(out, sd.StdInit(v))
	}

	return out
}

func (c *Context) tmplSetRoles(target interface{}, roleSlice interface{}) (string, error) {
	if c.IncreaseCheckGenericAPICall() {
		return "", ErrTooManyAPICalls
	}

	targetID := targetUserID(target)
	if targetID == 0 {
		return "", nil
	}

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

	roles := make([]string, 0, rSlice.Len())
	for i := 0; i < rSlice.Len(); i++ {
		switch v := rSlice.Index(i).Interface().(type) {
		case string:
			roles = append(roles, v)
		case int, int64:
			roles = append(roles, discordgo.StrID(reflect.ValueOf(v).Int()))
		default:
			return "", errors.New("Could not convert slice to string slice")
		}
	}

	err := common.BotSession.GuildMemberEdit(c.GS.ID, targetID, roles)
	if err != nil {
		return "", err
	}
	return "", nil
}
