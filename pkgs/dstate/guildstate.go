package dstate

import (
	"errors"
	"sync"
	"time"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var (
	ErrMemberNotFound  = errors.New("Member not found")
	ErrChannelNotFound = errors.New("Channel not found")
)

type GuildState struct {
	sync.RWMutex

	// ID is never mutated, so can be accessed without locking
	ID int64 `json:"id"`

	// The underlying guild, the members and channels fields shouldnt be used
	Guild *discordgo.Guild `json:"guild"`

	Members  map[int64]*MemberState  `json:"members"`
	Channels map[int64]*ChannelState `json:"channels" `

	MaxMessages          int           // Absolute max number of messages cached in a channel
	MaxMessageDuration   time.Duration // Max age of messages, if 0 ignored. (Only checks age whena new message is received on the channel)
	RemoveOfflineMembers bool

	userCache *Cache
}

// NewGuildstate creates a new guild state, it only uses the passed state to get settings from
// Pass nil to use default settings
func NewGuildState(guild *discordgo.Guild, state *State) *GuildState {
	gCop := new(discordgo.Guild)
	*gCop = *guild

	guildState := &GuildState{
		ID:       guild.ID,
		Guild:    gCop,
		Members:  make(map[int64]*MemberState),
		Channels: make(map[int64]*ChannelState),
	}

	if state != nil {
		guildState.MaxMessages = state.MaxChannelMessages
		guildState.MaxMessageDuration = state.MaxMessageAge
		guildState.RemoveOfflineMembers = state.RemoveOfflineMembers
		guildState.userCache = NewCache(state.cacheHits, state.cacheMiss)
	} else {
		guildState.userCache = NewCache(nil, nil)
	}

	for _, channel := range gCop.Channels {
		guildState.ChannelAddUpdate(false, channel)
	}

	if state != nil && state.TrackMembers {
		for _, member := range gCop.Members {
			guildState.MemberAddUpdate(false, member)
		}

		for _, presence := range gCop.Presences {
			guildState.PresenceAddUpdate(false, presence)
		}
	}

	gCop.Presences = nil
	gCop.Members = nil
	gCop.Emojis = nil
	gCop.Channels = nil

	return guildState
}

// InitCache intiializses the cache, assumes that g is locked
func (g *GuildState) InitCache(state *State) {
	g.userCache = NewCache(state.cacheHits, state.cacheMiss)
}

// StrID is the same as above but formats it in a string first
func (g *GuildState) StrID() string {
	return discordgo.StrID(g.ID)
}

// GuildUpdate updates the guild with new guild information
func (g *GuildState) GuildUpdate(lock bool, newGuild *discordgo.Guild) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	if newGuild.Roles == nil {
		newGuild.Roles = g.Guild.Roles
	}
	if newGuild.Emojis == nil {
		newGuild.Emojis = g.Guild.Emojis
	}
	if newGuild.VoiceStates == nil {
		newGuild.VoiceStates = g.Guild.VoiceStates
	}
	if newGuild.MemberCount == 0 {
		newGuild.MemberCount = g.Guild.MemberCount
	}

	// Create/update new channels
	*g.Guild = *newGuild
	for _, c := range newGuild.Channels {
		g.ChannelAddUpdate(false, c)
	}

	// Remove removed channels
	if newGuild.Channels != nil {
	OUTER:
		for _, checking := range g.Channels {
			for _, c := range newGuild.Channels {
				if c.ID == checking.ID {
					continue OUTER
				}
			}

			delete(g.Channels, checking.ID)
		}
	}
}

// LightCopy returns a light copy of the inner guild (no slices)
func (g *GuildState) LightCopy(lock bool) *discordgo.Guild {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	gCopy := new(discordgo.Guild)

	*gCopy = *g.Guild
	gCopy.Members = nil
	gCopy.Presences = nil
	gCopy.Channels = nil
	gCopy.VoiceStates = nil
	gCopy.Roles = nil

	return gCopy
}

// DeepCopy returns a deeper copy of the inner guild, with full copies of the specified slices
func (g *GuildState) DeepCopy(lock bool, copyRoles, copyVoiceStates, copyChannels bool) *discordgo.Guild {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	gCopy := new(discordgo.Guild)

	*gCopy = *g.Guild
	gCopy.Members = nil
	gCopy.Presences = nil
	gCopy.Channels = nil
	gCopy.VoiceStates = nil
	gCopy.Roles = nil

	if copyRoles {
		gCopy.Roles = make([]*discordgo.Role, len(g.Guild.Roles))
		copy(gCopy.Roles, g.Guild.Roles)
	}

	if copyVoiceStates {
		gCopy.VoiceStates = make([]*discordgo.VoiceState, len(g.Guild.VoiceStates))
		copy(gCopy.VoiceStates, g.Guild.VoiceStates)
	}

	if copyChannels {
		gCopy.Channels = make([]*discordgo.Channel, len(g.Channels))

		i := 0
		for _, v := range g.Channels {
			gCopy.Channels[i] = v.DGoCopy()
			i++
		}
	}

	return gCopy
}

// Member returns a the member from an id, or nil if not found
func (g *GuildState) Member(lock bool, id int64) *MemberState {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	return g.Members[id]
}

// MemberCopy returns a full copy of a MemberState, this can be used without locking
// Warning: modifying slices in the state (such as roles) causes race conditions, they're only safe to access
func (g *GuildState) MemberCopy(lock bool, id int64) *MemberState {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	m := g.Member(false, id)
	if m == nil {
		return nil
	}

	return m.Copy()
}

// MemberDGoCopy returns a full copy of a MemberState converted to a discordgo.Member struct
// Warning: modifying slices in the state (such as roles) causes race conditions, they're only safe to access
func (g *GuildState) MemberDGoCopy(lock bool, id int64) *discordgo.Member {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	m := g.Member(false, id)
	if m == nil {
		return nil
	}

	return m.DGoCopy()
}

// MemberAddUpdate adds or updates a member
func (g *GuildState) MemberAddUpdate(lock bool, newMember *discordgo.Member) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	if g.Members == nil {
		panic("nil members")
	}

	existing, ok := g.Members[newMember.User.ID]
	if ok {
		existing.UpdateMember(newMember)
	} else {
		ms := &MemberState{
			Guild: g,
			ID:    newMember.User.ID,
			Bot:   newMember.User.Bot,
		}

		ms.UpdateMember(newMember)
		g.Members[newMember.User.ID] = ms
	}
}

// MemberAdd adds a member to the GuildState
// It differs from addupdate in that it first increases the membercount and then calls MemberAddUpdate
// so it should only be called on the GuildMemberAdd event
func (g *GuildState) MemberAdd(lock bool, newMember *discordgo.Member) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	g.Guild.MemberCount++
	g.MemberAddUpdate(false, newMember)
}

// MemberRemove removes a member from the guildstate
// it also decrements membercount, so only call this on the GuildMemberRemove event
// If you wanna remove a member purely from the state, use StateRemoveMember
func (g *GuildState) MemberRemove(lock bool, id int64) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	g.Guild.MemberCount--
	delete(g.Members, id)
}

// StateRemoveMember removes a guildmember from the state and does NOT decrement member_count
func (g *GuildState) StateRemoveMember(lock bool, id int64) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	delete(g.Members, id)
}

// PresenceAddUpdate adds or updates a presence
func (g *GuildState) PresenceAddUpdate(lock bool, newPresence *discordgo.Presence) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	existing, ok := g.Members[newPresence.User.ID]
	if ok {
		existing.UpdatePresence(newPresence)
	} else {
		if newPresence.Status == discordgo.StatusOffline {
			// Don't bother if this is the case, most likely just removed from the server and the state would be very incomplete
			return
		}

		ms := &MemberState{
			Guild: g,
			ID:    newPresence.User.ID,
			Bot:   newPresence.User.Bot,
		}

		ms.UpdatePresence(newPresence)
		g.Members[newPresence.User.ID] = ms
	}
}

// ChannelCopy returns a copy of a channel
// Read actions are safe to do on the copy's slices, but not write actions
func (g *GuildState) ChannelCopy(lock bool, id int64) *ChannelState {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	c := g.Channel(false, id)
	if c == nil {
		return nil
	}

	return c.Copy(false)
}

// Channel retrieves a channelstate by id
func (g *GuildState) Channel(lock bool, id int64) *ChannelState {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	return g.Channels[id]
}

// ChannelAddUpdate adds or updates a channel in the guildstate
func (g *GuildState) ChannelAddUpdate(lock bool, newChannel *discordgo.Channel) *ChannelState {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	existing, ok := g.Channels[newChannel.ID]
	if ok {
		// Patch
		existing.Update(false, newChannel)
		return existing
	}

	state := NewChannelState(g, g, newChannel)
	g.Channels[newChannel.ID] = state

	return state
}

// ChannelRemove removes a channel from the GuildState
func (g *GuildState) ChannelRemove(lock bool, id int64) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}
	delete(g.Channels, id)
}

// Role returns a role by id, this is a strict copy
func (g *GuildState) RoleCopy(lock bool, id int64) *discordgo.Role {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	for _, role := range g.Guild.Roles {
		if role.ID == id {
			cop := *role
			return &cop
		}
	}

	return nil
}

// Role returns a role by id
func (g *GuildState) Role(lock bool, id int64) *discordgo.Role {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	for _, role := range g.Guild.Roles {
		if role.ID == id {
			return role
		}
	}

	return nil
}

func (g *GuildState) roleIndex(id int64) int {
	for i, v := range g.Guild.Roles {
		if v.ID == id {
			return i
		}
	}

	return -1
}

func (g *GuildState) RoleAddUpdate(lock bool, newRole *discordgo.Role) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	rCop := *newRole

	existingIndex := g.roleIndex(newRole.ID)

	if existingIndex == -1 {
		g.Guild.Roles = append(g.Guild.Roles, newRole)
	} else {
		g.Guild.Roles[existingIndex] = &rCop
	}
}

func (g *GuildState) RoleRemove(lock bool, id int64) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	index := g.roleIndex(id)
	if index == -1 {
		return
	}

	g.Guild.Roles = append(g.Guild.Roles[:index], g.Guild.Roles[index+1:]...)
}

func (g *GuildState) VoiceState(lock bool, userID int64) *discordgo.VoiceState {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	for _, v := range g.Guild.VoiceStates {
		if v.UserID == userID {
			return v
		}
	}

	return nil
}

func (g *GuildState) voiceStateIndex(userID int64) int {
	for i, v := range g.Guild.VoiceStates {
		if v.UserID == userID {
			return i
		}
	}

	return -1
}

func (g *GuildState) VoiceStateUpdate(lock bool, update *discordgo.VoiceState) {
	if lock {
		g.Lock()
		defer g.Unlock()
	}

	index := g.voiceStateIndex(update.UserID)
	if update.ChannelID == 0 {
		// left the channel
		if index == -1 {
			// was never in a channel?
			return
		}

		g.Guild.VoiceStates = append(g.Guild.VoiceStates[:index], g.Guild.VoiceStates[index+1:]...)
		return
	}

	vsCopy := new(discordgo.VoiceState)
	*vsCopy = *update

	if index != -1 {
		g.Guild.VoiceStates[index] = vsCopy
	} else {
		g.Guild.VoiceStates = append(g.Guild.VoiceStates, vsCopy)
	}
}

// Calculates the permissions for a member.
// https://support.discordapp.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func (g *GuildState) MemberPermissions(lock bool, channelID int64, memberID int64) (apermissions int64, err error) {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	if memberID == g.Guild.OwnerID {
		return discordgo.PermissionAll, nil
	}

	mState := g.Member(false, memberID)
	if mState == nil {
		return 0, ErrMemberNotFound
	}

	return g.MemberPermissionsMS(false, channelID, mState)
}

// Calculates the permissions for a member.
// https://support.discordapp.com/hc/en-us/articles/206141927-How-is-the-permission-hierarchy-structured-
func (g *GuildState) MemberPermissionsMS(lock bool, channelID int64, mState *MemberState) (apermissions int64, err error) {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	if mState.ID == g.Guild.OwnerID {
		return discordgo.PermissionAll, nil
	}

	for _, role := range g.Guild.Roles {
		if role.ID == g.Guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range g.Guild.Roles {
		for _, roleID := range mState.Roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	// Administrator bypasses channel overrides
	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAll
		return
	}

	cState := g.Channel(false, channelID)
	if cState == nil {
		err = ErrChannelNotFound
		return
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range cState.PermissionOverwrites {
		if g.Guild.ID == overwrite.ID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	var denies, allows int64

	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range cState.PermissionOverwrites {
		for _, roleID := range mState.Roles {
			if overwrite.Type == discordgo.PermissionOverwriteTypeRole && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	apermissions &= ^denies
	apermissions |= allows

	for _, overwrite := range cState.PermissionOverwrites {
		if overwrite.Type == discordgo.PermissionOverwriteTypeMember && overwrite.ID == mState.ID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAllChannel
	}

	return
}

func (g *GuildState) runGC(cacheExpirey time.Duration, offlineMembers bool, maxMessages int, maxMessageDuration time.Duration) (membersEvicted int, cacheN int, messagesRemoved int) {
	if g.userCache != nil {
		cacheN = g.userCache.EvictOldKeys(time.Now().Add(-cacheExpirey))
	}

	g.Lock()
	defer g.Unlock()
	if offlineMembers {
		for k, v := range g.Members {
			if !v.PresenceSet || v.PresenceStatus == StatusOffline {
				delete(g.Members, k)
				membersEvicted++
			}
		}
	}

	t := time.Now()
	for _, v := range g.Channels {
		messagesRemoved += v.runGC(t, maxMessageDuration, maxMessages)
	}

	return
}

// UserCacheGet retrieves a item from the cache
// Safe to call without locking GuildState as there's another lock managed by the cache internally
func (g *GuildState) UserCacheGet(key interface{}) interface{} {
	if g.userCache == nil {
		return nil
	}

	return g.userCache.Get(key)
}

// UserCacheSet stores an item in the cache
// Safe to call without locking GuildState as there's another lock managed by the cache internally
func (g *GuildState) UserCacheSet(key interface{}, value interface{}) {
	if g.userCache == nil {
		return
	}

	g.userCache.Set(key, value)
}

// UserCacheDel deletes an item from the cache
// Safe to call without locking GuildState as there's another lock managed by the cache internally
func (g *GuildState) UserCacheDel(key interface{}) {
	if g.userCache == nil {
		return // nothing to delete
	}

	g.userCache.Del(key)
}

// UserCacheFetch either retrieves an existing item from the cache or fetches one from the provided CacheFetchFunc
// Safe to call without locking GuildState as there's another lock managed by the cache internally
func (g *GuildState) UserCacheFetch(key interface{}, fetchFunc CacheFetchFunc) (interface{}, error) {
	if g.userCache == nil {
		return nil, errors.New("No cache")
	}

	return g.userCache.Fetch(key, fetchFunc)
}

// UserCacheDellAllKeysType deletes all keys of the provided type
func (g *GuildState) UserCacheDellAllKeysType(key interface{}) {
	if g.userCache == nil {
		return
	}

	g.userCache.DelAllKeysType(key)
}

// IsAvailable returns wether the guild is available or not (guild outages or starting up)
func (g *GuildState) IsAvailable(lock bool) bool {
	if lock {
		g.RLock()
		defer g.RUnlock()
	}

	return !g.Guild.Unavailable
}
