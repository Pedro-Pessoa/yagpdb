package dstate

import (
	"testing"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

func createTestGS() *GuildState {
	state := NewState()
	return NewGuildState(createTestGuild(1, 2), state)
}

func TestGuildRoles(t *testing.T) {
	gs := createTestGS()

	// Test adding roles
	preLen := len(gs.Guild.Roles)
	gs.RoleAddUpdate(true, &discordgo.Role{ID: 50, Name: "t1"})

	nowLen := len(gs.Guild.Roles)
	if nowLen != preLen+1 {
		t.Errorf("role not added: %#v", gs.Guild.Roles)
		return
	}

	gs.RoleAddUpdate(true, &discordgo.Role{ID: 51, Name: "t2"})
	nowLen = len(gs.Guild.Roles)
	if nowLen != preLen+2 {
		t.Errorf("role2 not added: %#v", gs.Guild.Roles)
		return
	}

	// make sure they're correct
	r1 := gs.RoleCopy(true, 50)
	if r1.ID != 50 {
		t.Errorf("add: r1 incorrect: f:%#v", r1)
		return
	}

	r2 := gs.RoleCopy(true, 51)
	if r2.ID != 51 {
		t.Errorf("add: r2 incorrect: f:%#v", r2)
		return
	}

	// test updating roles, r1
	gs.RoleAddUpdate(true, &discordgo.Role{ID: 50, Name: "new name"})

	r1 = gs.RoleCopy(true, 50)
	if r1.ID != 50 || r1.Name != "new name" {
		t.Errorf("change p1: r1 incorrect: f:%#v", r1)
		return
	}

	r2 = gs.RoleCopy(true, 51)
	if r2.ID != 51 || r2.Name != "t2" {
		t.Errorf("change p1: r2 incorrect: f:%#v", r2)
		return
	}

	// test updating roles, r2
	gs.RoleAddUpdate(true, &discordgo.Role{ID: 51, Name: "new name2"})

	r1 = gs.RoleCopy(true, 50)
	if r1.ID != 50 || r1.Name != "new name" {
		t.Errorf("change p2: r1 incorrect: f:%#v", r1)
		return
	}

	r2 = gs.RoleCopy(true, 51)
	if r2.ID != 51 || r2.Name != "new name2" {
		t.Errorf("change p2: r2 incorrect: f:%#v", r2)
		return
	}

	// test removing roles p1
	gs.RoleRemove(true, 50)

	r1 = gs.RoleCopy(true, 50)
	if r1 != nil {
		t.Errorf("del p1: r1 incorrect: f:%#v", r1)
		return
	}

	r2 = gs.RoleCopy(true, 51)
	if r2.ID != 51 || r2.Name != "new name2" {
		t.Errorf("del p1: r2 incorrect: f:%#v", r2)
		return
	}

	// test removing roles p2
	gs.RoleRemove(true, 51)

	r1 = gs.RoleCopy(true, 50)
	if r1 != nil {
		t.Errorf("del p1: r1 incorrect: f:%#v", r1)
		return
	}

	r2 = gs.RoleCopy(true, 51)
	if r2 != nil {
		t.Errorf("del p1: r2 incorrect: f:%#v", r2)
		return
	}
}

func TestGuildVoiceStates(t *testing.T) {
	gs := createTestGS()

	// Test adding roles
	preLen := len(gs.Guild.VoiceStates)
	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 50, ChannelID: 2})

	nowLen := len(gs.Guild.VoiceStates)
	if nowLen != preLen+1 {
		t.Errorf("vs not added: %#v", gs.Guild.VoiceStates)
		return
	}

	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 51, ChannelID: 2})
	nowLen = len(gs.Guild.VoiceStates)
	if nowLen != preLen+2 {
		t.Errorf("vs2 not added: %#v", gs.Guild.VoiceStates)
		return
	}

	// make sure they're correct
	vs1 := gs.VoiceState(true, 50)
	if vs1.UserID != 50 {
		t.Errorf("add: vs1 incorrect: %#v", vs1)
		return
	}

	vs2 := gs.VoiceState(true, 51)
	if vs2.UserID != 51 {
		t.Errorf("add: vs2 incorrect: %#v", vs2)
		return
	}

	// test updating states, vs1
	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 50, ChannelID: 3})

	vs1 = gs.VoiceState(true, 50)
	if vs1.UserID != 50 || vs1.ChannelID != 3 {
		t.Errorf("change p1: vs1 incorrect: %#v", vs1)
		return
	}

	vs2 = gs.VoiceState(true, 51)
	if vs2.UserID != 51 || vs2.ChannelID != 2 {
		t.Errorf("change p1: vs2 incorrect: %#v", vs2)
		return
	}

	// test updating states, vs2
	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 51, ChannelID: 4})

	vs1 = gs.VoiceState(true, 50)
	if vs1.UserID != 50 || vs1.ChannelID != 3 {
		t.Errorf("change p2: vs1 incorrect: %#v", vs1)
		return
	}

	vs2 = gs.VoiceState(true, 51)
	if vs2.UserID != 51 || vs2.ChannelID != 4 {
		t.Errorf("change p2: vs2 incorrect: %#v", vs2)
		return
	}

	// test removing states p1
	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 50, ChannelID: 0})

	vs1 = gs.VoiceState(true, 50)
	if vs1 != nil {
		t.Errorf("del p1: vs1 incorrect: %#v", vs1)
		return
	}

	vs2 = gs.VoiceState(true, 51)
	if vs2.UserID != 51 || vs2.ChannelID != 4 {
		t.Errorf("del p1: vs2 incorrect: %#v", vs2)
		return
	}

	// test removing states p2
	gs.VoiceStateUpdate(true, &discordgo.VoiceState{UserID: 51, ChannelID: 0})

	vs1 = gs.VoiceState(true, 50)
	if vs1 != nil {
		t.Errorf("del p1: vs1 incorrect: %#v", vs1)
		return
	}

	vs2 = gs.VoiceState(true, 51)
	if vs2 != nil {
		t.Errorf("del p1: vs2 incorrect: %#v", vs2)
		return
	}
}
