package moderation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

type ModlogAction struct {
	Prefix string
	Emoji  string
	Color  int
	Name   string
	Footer string
}

func (m ModlogAction) String() string {
	str := m.Emoji + m.Prefix
	if m.Footer != "" {
		str += " (" + m.Footer + ")"
	}

	return str
}

var (
	MAMute       = ModlogAction{Prefix: "Silenciado:", Emoji: "🔇", Color: 0x57728e, Name: "🔇Silenciado!"}
	MAUnmute     = ModlogAction{Prefix: "Desilenciado:", Emoji: "🔊", Color: 0x62c65f, Name: "🔊Desilenciado!"}
	MAKick       = ModlogAction{Prefix: "Quicado:", Emoji: "👢", Color: 0xf2a013, Name: "👢Quicado!"}
	MABanned     = ModlogAction{Prefix: "Banido:", Emoji: "🔨", Color: 0xd64848, Name: "🔨Banido!"}
	MAUnbanned   = ModlogAction{Prefix: "Desbanido:", Emoji: "🔓", Color: 0x62c65f}
	MAWarned     = ModlogAction{Prefix: "Avisado:", Emoji: "⚠", Color: 0xfca253, Name: "⚠Avisado!"}
	MASlowmode   = ModlogAction{Prefix: "Freiou", Emoji: "🐌", Color: 0xff9a00}
	MARemoveSlow = ModlogAction{Prefix: "Desfreiou", Emoji: "💨", Color: 0x43ff00}
	MAGiveRole   = ModlogAction{Prefix: "", Emoji: "✅", Color: 0x53fcf9}
	MARemoveRole = ModlogAction{Prefix: "", Emoji: "🛑", Color: 0x53fcf9}
	MALock       = ModlogAction{Prefix: "trancado", Emoji: "🔒", Color: 0xEE00EE}
	MAUnlock     = ModlogAction{Prefix: "destrancado", Emoji: "🔓", Color: 0x718AED}
)

func CreateModlogEmbed(config *Config, author *discordgo.User, action ModlogAction, target interface{}, reason, logLink string) error {
	channelID := config.IntActionChannel()
	config.GetGuildID()
	if channelID == 0 {
		return nil
	}

	emptyAuthor := false
	if author == nil {
		emptyAuthor = true
		author = &discordgo.User{
			ID:            0,
			Username:      "Unknown",
			Discriminator: "????",
		}
	}

	name := ""
	discriminator := ""
	var id int64
	switch t := target.(type) {
	case *discordgo.User:
		name = t.Username
		discriminator = "#" + t.Discriminator
		id = t.ID
	case *discordgo.Role:
		name = t.Name
		id = t.ID
	case *discordgo.Channel:
		name = t.Name
		id = t.ID
	}

	if reason == "" {
		reason = "(sem motivo especificado)"
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s#%s (ID %d)", author.Username, author.Discriminator, author.ID),
			IconURL: discordgo.EndpointUserAvatar(author.ID, author.Avatar),
		},
		Color: action.Color,
		Description: fmt.Sprintf("**%s%s %s**%s *(ID %d)*\n📄**Motivo:** %s",
			action.Emoji, action.Prefix, name, discriminator, id, reason),
	}

	if discriminator != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: discordgo.EndpointUserAvatar(id, target.(*discordgo.User).Avatar),
		}
	}

	if logLink != "" {
		embed.Description += " ([Logs](" + logLink + "))"
	}

	if action.Footer != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: action.Footer,
		}
	}

	m, err := common.BotSession.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		if common.IsDiscordErr(err, discordgo.ErrCodeMissingAccess, discordgo.ErrCodeMissingPermissions, discordgo.ErrCodeUnknownChannel) {
			// disable the modlog
			config.ActionChannel = ""
			_ = config.Save(config.GetGuildID())
			return nil
		}
		return err
	}

	if emptyAuthor {
		placeholder := fmt.Sprintf("Coloque um autor usando **'reason %d o-motivo-aqui`**", m.ID)
		updateEmbedReason(nil, placeholder, embed)
		_, err = common.BotSession.ChannelMessageEditEmbed(channelID, m.ID, embed)
	}

	return err
}

var (
	logsRegex = regexp.MustCompile(`\(\[Logs\]\(.*\)\)`)
)

func updateEmbedReason(author *discordgo.User, reason string, embed *discordgo.MessageEmbed) {
	const checkStr = "📄**Motivo:**"

	index := strings.Index(embed.Description, checkStr)
	withoutReason := embed.Description[:index+len(checkStr)]

	logsLink := logsRegex.FindString(embed.Description)
	if logsLink != "" {
		logsLink = " " + logsLink
	}

	embed.Description = withoutReason + " " + reason + logsLink

	if author != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s#%s (ID %d)", author.Username, author.Discriminator, author.ID),
			IconURL: discordgo.EndpointUserAvatar(author.ID, author.Avatar),
		}
	}
}
