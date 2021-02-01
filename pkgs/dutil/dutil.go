package dutil

// Package dutil provides general discordgo utilities that i find to be reusing across my discord projects

import (
	"regexp"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

// GetAllGuildMembers Returns all guild members in a guild
// It will make `number of members`/1000 requests to the api
func GetAllGuildMembers(session *discordgo.Session, guilID int64) ([]*discordgo.Member, error) {
	var after int64
	members := make([]*discordgo.Member, 0)

	for {
		resp, err := session.GuildMembers(guilID, after, 1000)
		if err != nil {
			return nil, err
		}
		members = append(members, resp...)

		if len(resp) < 1000 {
			break // Reached the end
		}

		after = members[len(members)-1].User.ID
	}
	return members, nil
}

// IsRoleAbove returns wether role a is above b, checking positions first, and if they're the same
// (both being 1, new roles always have 1 as position)
// then it checjs by lower id
func IsRoleAbove(a, b *discordgo.Role) bool {
	if a.Position != b.Position {
		return a.Position > b.Position
	}

	if a.ID == b.ID {
		return false
	}

	return a.ID < b.ID
}

// Channels are a collection of Channels
type Channels []*discordgo.Channel

func (r Channels) Len() int {
	return len(r)
}

func (r Channels) Less(i, j int) bool {
	return r[i].Position < r[j].Position
}

func (r Channels) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Roles []*discordgo.Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return IsRoleAbove(r[i], r[j])
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

const zeroWidthSpace = "​"

var (
	everyoneReplacer    = strings.NewReplacer("@everyone", "@"+zeroWidthSpace+"everyone")
	hereReplacer        = strings.NewReplacer("@here", "@"+zeroWidthSpace+"here")
	patternRoleMentions = regexp.MustCompile("<@&[0-9]*>")
)

// EscapeEveryoneMention Escapes an everyone mention, adding a zero width space between the '@' and rest
func EscapeEveryoneMention(in string) string {
	s := everyoneReplacer.Replace(in)
	s = hereReplacer.Replace(s)
	return s
}

// EscapeSpecialMentions Escapes an everyone mention, adding a zero width space between the '@' and rest, also escapes role mentions
func EscapeSpecialMentions(s string) string {
	s = EscapeEveryoneMention(s)

	return patternRoleMentions.ReplaceAllStringFunc(s, func(x string) string {
		if len(x) < 4 {
			return x
		}

		// Not allowed
		return x[:2] + zeroWidthSpace + x[2:]
	})
}
