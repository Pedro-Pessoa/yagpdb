package discordgo

// Webhook stores the data for a webhook.
type Webhook struct {
	ID        int64       `json:"id,string"`
	Type      WebhookType `json:"type"`
	GuildID   int64       `json:"guild_id,string"`
	ChannelID int64       `json:"channel_id,string"`
	User      *User       `json:"user"`
	Name      string      `json:"name"`
	Avatar    string      `json:"avatar"`
	Token     string      `json:"token"`

	// ApplicationID is the bot/OAuth2 application that created this webhook
	ApplicationID int64 `json:"application_id,string,omitempty"`
}

// WebhookType is the type of Webhook (see WebhookType* consts) in the Webhook struct
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type WebhookType int

// Valid WebhookType values
const (
	WebhookTypeIncoming WebhookType = iota
	WebhookTypeChannelFollower
)

// WebhookParams is a struct for webhook params, used in the WebhookExecute command.
type WebhookParams struct {
	Content         string                  `json:"content,omitempty"`
	Username        string                  `json:"username,omitempty"`
	AvatarURL       string                  `json:"avatar_url,omitempty"`
	TTS             bool                    `json:"tts,omitempty"`
	Files           []*File                 `json:"-"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}

// WebhookEdit stores data for editing of a webhook message.
type WebhookEdit struct {
	Content         string                  `json:"content,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}
