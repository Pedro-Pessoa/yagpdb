package simpleembed

import (
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	CmdCategory:         commands.CategoryFun,
	Name:                "SimpleEmbed",
	Aliases:             []string{"se"},
	Description:         "A more simpler version of CustomEmbed, controlled completely using switches.",
	RequireDiscordPerms: []int64{discordgo.PermissionManageMessages},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "channel", Help: "Optional channel to send in", Type: dcmd.Channel},
		{Switch: "content", Help: "Text content for the message", Type: dcmd.String, Default: ""},

		{Switch: "title", Type: dcmd.String, Default: ""},
		{Switch: "desc", Type: dcmd.String, Help: "Text in the 'description' field", Default: ""},
		{Switch: "color", Help: "Either hex code or name", Type: dcmd.String, Default: ""},
		{Switch: "url", Help: "Url of this embed", Type: dcmd.String, Default: ""},
		{Switch: "thumbnail", Help: "Url to a thumbnail", Type: dcmd.String, Default: ""},
		{Switch: "image", Help: "Url to an image", Type: dcmd.String, Default: ""},

		{Switch: "author", Help: "The text in the 'author' field", Type: dcmd.String, Default: ""},
		{Switch: "authoricon", Help: "Url to a icon for the 'author' field", Type: dcmd.String, Default: ""},
		{Switch: "authorurl", Help: "Url of the 'author' field", Type: dcmd.String, Default: ""},

		{Switch: "footer", Help: "Text content for the footer", Type: dcmd.String, Default: ""},
		{Switch: "footericon", Help: "Url to a icon for the 'footer' field", Type: dcmd.String, Default: ""},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		content := data.Switch("content").Str()
		embed := &discordgo.MessageEmbed{
			Title:       data.Switch("title").Str(),
			Description: data.Switch("desc").Str(),
			URL:         data.Switch("url").Str(),
		}

		if color := data.Switch("color").Str(); color != "" {
			parsedColor, ok := util.ParseColor(color)
			if !ok {
				return "Unknown color: " + color + ", can be either hex color code or name for a known color", nil
			}

			embed.Color = parsedColor
		}

		if author := data.Switch("author").Str(); author != "" {
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    author,
				IconURL: data.Switch("authoricon").Str(),
				URL:     data.Switch("authorurl").Str(),
			}
		}

		if thumbnail := data.Switch("thumbnail").Str(); thumbnail != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: thumbnail,
			}
		}

		if image := data.Switch("image").Str(); image != "" {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: image,
			}
		}

		footer := data.Switch("footer").Str()
		footerIcon := data.Switch("footericon").Str()
		if footer != "" || footerIcon != "" {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text:    footer,
				IconURL: footerIcon,
			}
		}

		cID := data.Msg.ChannelID
		c := data.Switch("channel")
		if c.Value != nil {
			cID = c.Value.(*dstate.ChannelState).ID

			hasPerms, err := bot.AdminOrPermMS(cID, data.MS, discordgo.PermissionSendMessages|discordgo.PermissionViewChannel)
			if err != nil {
				return "Failed checking permissions, please try again or join the support server.", err
			}

			if !hasPerms {
				return "You do not have permissions to send messages there", nil
			}
		}

		messageSend := &discordgo.MessageSend{
			Content:         content,
			Embed:           embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		}

		_, err := common.BotSession.ChannelMessageSendComplex(cID, messageSend)
		if err != nil {
			return err, err
		}

		if cID != data.Msg.ChannelID {
			return "Done", nil
		}

		return nil, nil
	},
}
