package dstate

import (
	"strconv"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

// MessageState represents the state of a message
type MessageState struct {

	// The ID of the message.
	ID int64 `json:"id,string"`

	// The ID of the channel in which the message was sent.
	ChannelID int64 `json:"channel_id,string"`

	// The ID of the guild in which the message was sent.
	GuildID int64 `json:"guild_id,string,omitempty"`

	// The content of the message.
	Content string `json:"content"`

	// The roles mentioned in the message.
	MentionRoles []int64 `json:"mention_roles"`

	// Whether the message is text-to-speech.
	Tts bool `json:"tts"`

	// Whether the message mentions everyone.
	MentionEveryone bool `json:"mention_everyone"`

	// The author of the message. This is not guaranteed to be a
	// valid user (webhook-sent messages do not possess a full author).
	Author discordgo.User `json:"author"`

	// A list of attachments present in the message.
	Attachments []*discordgo.MessageAttachment `json:"attachments"`

	// A list of embeds present in the message. Multiple
	// embeds can currently only be sent by webhooks.
	Embeds []*discordgo.MessageEmbed `json:"embeds"`

	// A list of users mentioned in the message.
	Mentions []*discordgo.User `json:"mentions"`

	// The type of the message.
	Type discordgo.MessageType `json:"type"`

	WebhookID int64 `json:"webhook_id,string"`

	///////////////////////////////////////
	// NON STANDARD MESSAGE FIELDS BELOW //
	///////////////////////////////////////

	// Set it the message was deleted
	Deleted bool

	// The parsed times below are cached because parsing all messages
	// timestamps in state everytime a new one came in would be stupid
	ParsedCreated time.Time
	ParsedEdited  time.Time
}

func MessageStateFromMessage(msg *discordgo.Message) *MessageState {
	var author discordgo.User
	if msg.Author != nil {
		author = *msg.Author
	}

	ms := &MessageState{
		ID:              msg.ID,
		ChannelID:       msg.ChannelID,
		GuildID:         msg.GuildID,
		Content:         msg.Content,
		MentionRoles:    msg.MentionRoles,
		Tts:             msg.TTS,
		MentionEveryone: msg.MentionEveryone,
		Author:          author,
		Attachments:     msg.Attachments,
		Embeds:          msg.Embeds,
		Mentions:        msg.Mentions,
		Type:            msg.Type,
		WebhookID:       msg.WebhookID,
	}

	ms.parseTimes(msg.Timestamp, msg.EditedTimestamp)
	return ms
}

// ParseTimes parses the created and edited timestamps
func (m *MessageState) parseTimes(created, edited discordgo.Timestamp) {
	// The discord api is stopid, and edits can come before creates
	// Can also be handled before even if received in order cause of goroutines and scheduling
	if created != "" {
		parsedC, _ := created.Parse()
		m.ParsedCreated = parsedC
	}

	if edited != "" {
		parsedE, _ := edited.Parse()
		m.ParsedEdited = parsedE
	}
}

// Copy returns a copy of the message, that can be used without further locking the owner
func (m *MessageState) Copy() *MessageState {
	cop := *m
	return &cop
}

func (m *MessageState) Update(msg *discordgo.Message) {
	// Patch the m message
	if msg.Content != "" {
		m.Content = msg.Content
	}
	if msg.Mentions != nil {
		m.Mentions = msg.Mentions
	}
	if msg.Embeds != nil {
		m.Embeds = msg.Embeds
	}
	if msg.Attachments != nil {
		m.Attachments = msg.Attachments
	}
	if msg.Author != nil {
		m.Author = *msg.Author
	}
	if msg.MentionRoles != nil {
		m.MentionRoles = msg.MentionRoles
	}

	m.parseTimes(msg.Timestamp, msg.EditedTimestamp)
}

func IsPrivate(t discordgo.ChannelType) bool {
	return t == discordgo.ChannelTypeGroupDM || t == discordgo.ChannelTypeDM
}

// ContentWithMentionsReplaced will replace all @<id> mentions with the
// username of the mention.
func (m *MessageState) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+strconv.FormatInt(user.ID, 10)+">", "@"+user.Username,
			"<@!"+strconv.FormatInt(user.ID, 10)+">", "@"+user.Username,
		).Replace(content)
	}
	return
}
