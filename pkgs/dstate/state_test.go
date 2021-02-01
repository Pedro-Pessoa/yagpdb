package dstate

import (
	"testing"
	"time"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var testState *State

func init() {
	testState = NewState()
	testGuild := createTestGuild(0, 1)
	testState.GuildCreate(false, testGuild)
}

func createTestGuild(gID, cID int64) *discordgo.Guild {
	return &discordgo.Guild{
		ID:   gID,
		Name: "Test Guild",
		Channels: []*discordgo.Channel{
			{ID: cID, Name: "Test Channel"},
		},
	}
}

func createTestMessage(mID, cID int64, content string) *discordgo.Message {
	return &discordgo.Message{ID: mID, Content: content, ChannelID: cID}
}

func TestGuildCreate(t *testing.T) {
	g := createTestGuild(100, 200)
	s := NewState()
	s.GuildCreate(true, g)

	// Check if guild got added
	gs := s.Guild(true, 100)
	if gs == nil {
		t.Fatal("GuildState is nil")
	}

	// Check if channel got added
	cs := s.Channel(true, 200)
	if cs == nil {
		t.Fatal("ChannelState is nil in global map")
	}

	cs = gs.Channel(true, 200)
	if cs == nil {
		t.Fatal("ChannelState is nil in guildstate map")
	}
}

func TestSecondReady(t *testing.T) {
	r := &discordgo.Ready{
		Guilds: []*discordgo.Guild{
			{
				ID:          1,
				Name:        "G 1",
				Unavailable: true,
			},
		},
	}

	s := NewState()
	s.HandleReady(r)

	loadGuilds := func() {
		for _, v := range r.Guilds {
			g := &discordgo.Guild{
				ID:   v.ID,
				Name: v.Name,
				Channels: []*discordgo.Channel{
					{
						ID:   v.ID + 1,
						Name: "C " + discordgo.StrID(v.ID+1),
					},
					{
						ID:   v.ID + 2,
						Name: "C " + discordgo.StrID(v.ID+2),
					},
					{
						ID:   v.ID + 3,
						Name: "C " + discordgo.StrID(v.ID+3),
					},
				},
			}
			s.GuildCreate(true, g)
		}
	}

	loadGuilds()

	doChecks := func(prefix string) {
		gs := s.Guild(true, 1)
		if gs == nil {
			t.Fatal(prefix + " GuildState is nil")
		}

		csLocal := gs.Channel(true, 2)
		if csLocal == nil {
			t.Fatal(prefix + " csLocal == nil")
		}

		csGlobal := s.Channel(true, 2)
		if csGlobal == nil {
			t.Fatal(prefix + " csGlobal == nil")
		}
	}
	doChecks("Initial:")

	s.HandleReady(r)

	doChecks("SecondReady:")

	loadGuilds()

	doChecks("SecondLoad:")
}

func TestGuildDelete(t *testing.T) {
	s := NewState()
	g := createTestGuild(100, 200)
	s.GuildCreate(true, g)

	s.GuildRemove(100)

	// Check if guild got removed
	gs := s.Guild(true, 100)
	if gs != nil {
		t.Fatal("GuildState is not nil")
	}

	// Check if channel got removed
	cs := s.Channel(true, 200)
	if cs != nil {
		t.Fatal("ChannelState is not nil in global map")
	}
}

func TestMessageCreate(t *testing.T) {
	s := NewState()
	s.MaxChannelMessages = 100
	g := createTestGuild(100, 200)
	s.GuildCreate(true, g)

	msgEvt1 := &discordgo.MessageCreate{
		Message: createTestMessage(300, 200, "Hello there buddy"),
	}
	msgEvt2 := &discordgo.MessageCreate{
		Message: createTestMessage(301, 200, "Hello there buddy"),
	}

	cs := s.Channel(true, 200)
	if cs == nil {
		t.Fatal("ChannelState is nil")
	}

	s.HandleEvent(nil, msgEvt1)
	s.HandleEvent(nil, msgEvt2)

	if len(cs.Messages) != 2 {
		t.Fatal("Length of messages not 4:", cs.Messages)
	}

	for i := 0; i < 150; i++ {
		s.HandleEvent(nil, &discordgo.MessageCreate{
			Message: createTestMessage(302+int64(i), 200, "HHeyyy"),
		})
	}

	cs.runGC(time.Now(), 0, 100)

	if len(cs.Messages) != 100 {
		t.Fatal("Length of messages not 100:", len(cs.Messages))
	}
}

func BenchmarkMessageCreate(b *testing.B) {
	s := NewState()
	s.MaxChannelMessages = 100

	g := createTestGuild(100, 200)
	s.GuildCreate(true, g)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgEvt := &discordgo.MessageCreate{
			Message: createTestMessage(300+int64(i), 200, "Hello there buddy"),
		}

		s.HandleEvent(nil, msgEvt)
	}
}

func BenchmarkMessageCreateParallel(b *testing.B) {
	s := NewState()
	s.MaxChannelMessages = 100

	g := createTestGuild(100, 200)
	s.GuildCreate(true, g)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			msgEvt := &discordgo.MessageCreate{
				Message: createTestMessage(300+int64(i), 200, "Hello there buddy"),
			}
			s.HandleEvent(nil, msgEvt)
			i++
		}
	})
}

func BenchmarkMessageCreateParallelMultiGuild100(b *testing.B) {
	s := NewState()
	s.MaxChannelMessages = 100

	for i := int64(0); i < 100; i++ {
		g := createTestGuild(100+i, 200+i)
		s.GuildCreate(true, g)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			msgEvt := &discordgo.MessageCreate{
				Message: createTestMessage(300+int64(i), 200+int64(i%100), "Hello there buddy"),
			}
			s.HandleEvent(nil, msgEvt)
			i++
		}
	})
}

// func BenchmarkDGOStateMessageCreatePalellMultiGuild100(b *testing.B) {
// 	s := discordgo.NewState()
// 	s.MaxMessageCount = 100

// 	for i := 0; i < 100; i++ {
// 		g := &discordgo.Guild{
// 			ID: fmt.Sprintf("g%d", i),
// 			Channels: []*discordgo.Channel{
// 				&discordgo.Channel{ID: fmt.Sprint(i), Name: fmt.Sprint(i)},
// 			},
// 		}
// 		s.OnInterface(nil, &discordgo.GuildCreate{g})
// 	}

// 	idMap := genStringIdMap(b.N)

// 	b.ResetTimer()

// 	b.RunParallel(func(pb *testing.PB) {
// 		i := 0
// 		for pb.Next() {
// 			msgEvt := &discordgo.MessageCreate{
// 				Message: createTestMessage(idMap[i], idMap[i%100], "Hello there buddy"),
// 			}
// 			s.OnInterface(nil, msgEvt)
// 			i++
// 		}
// 	})
// }
