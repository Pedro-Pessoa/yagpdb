package serverinfo

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	prfx "github.com/Pedro-Pessoa/tidbot/common/prefix"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dstate"
	"github.com/Pedro-Pessoa/tidbot/premium"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "serverinfo",
	Aliases:     []string{"sinfo"},
	Description: "Shows some informations about the server",
	Cooldown:    5,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		embed := embedCreator(data, nil, false)
		return embed, nil
	},
}

var AdminCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryGeneral,
	Name:                 "serverinfoadm",
	Aliases:              []string{"sinfoadm"},
	Description:          "Get targeted server infos",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "invite", Name: "Create Invite?"},
	},
	RunFunc: util.RequireBotAdmin(func(data *dcmd.Data) (interface{}, error) {
		if data.Args[0].Value == nil {
			return embedCreator(data, nil, false), nil
		}

		gs := bot.State.Guild(false, data.Args[0].Int64())
		if gs == nil {
			return "Guild is not in state", nil
		}

		var createInvite bool
		if data.Switches["invite"].Value != nil && data.Switches["invite"].Value.(bool) {
			createInvite = true
		}

		return embedCreator(nil, gs, createInvite), nil
	}),
}

func embedCreator(data *dcmd.Data, customGuild *dstate.GuildState, createInvite bool) *discordgo.MessageEmbed {
	var gs *dstate.GuildState
	switch {
	case customGuild != nil:
		gs = customGuild
	case data != nil:
		gs = data.GS
	default:
		return nil
	}

	guild := gs.DeepCopy(false, true, true, true, true)
	guild.Members = make([]*discordgo.Member, len(gs.Members))
	i := 0
	for _, m := range gs.Members {
		guild.Members[i] = m.DGoCopy()
		i++
	}

	title := fmt.Sprintf("Info for %s", guild.Name)
	ownerID := guild.OwnerID
	description := guild.Description

	var features, feat string

	featMap := map[string]string{
		"INVITE_SPLASH":                    "Invite background",
		"VIP_REGIONS":                      "Upgraded voice bitrate",
		"VANITY_URL":                       "Vanity URL",
		"VERIFIED":                         "Verified",
		"PARTNERED":                        "Partnered",
		"COMMUNITY":                        "Community server",
		"COMMERCE":                         "Commerce",
		"NEWS":                             "News channel",
		"DISCOVERABLE":                     "Discoverable",
		"FEATURABLE":                       "Featurable",
		"ANIMATED_ICON":                    "Animated icon",
		"BANNER":                           "Banner image",
		"WELCOME_SCREEN_ENABLED":           "Welcome screen",
		"MEMBER_VERIFICATION_GATE_ENABLED": "Member screening",
		"PREVIEW_ENABLED":                  "Previewable",
	}

	for _, f := range guild.Features {
		features += "<:green:806413632912490507> " + featMap[f] + "\n"
	}

	if features == "" {
		feat = "<:red:806413641078407188> No Features"
	} else {
		feat = features
	}

	var channelOutput string
	var textCount, textLockedCount, categoryCount, newsCount, voiceCount, voiceLockedCount, storeCount, total int
	var everyoneID int64
	var inviteSet, inviteErr bool
	var invite *discordgo.Invite
	var err error

	for _, r := range guild.Roles {
		if r.Name == "@everyone" {
			everyoneID = r.ID
			break
		}
	}

	for _, c := range guild.Channels {
		total++
		switch c.Type {
		case discordgo.ChannelTypeGuildCategory:
			categoryCount++
		case discordgo.ChannelTypeGuildNews:
			newsCount++
		case discordgo.ChannelTypeGuildStore:
			storeCount++
		case discordgo.ChannelTypeGuildText:
			textCount++
			for _, ow := range c.PermissionOverwrites {
				if ow.Type == discordgo.PermissionOverwriteTypeRole && ow.ID == everyoneID && (ow.Deny&discordgo.PermissionViewChannel) == discordgo.PermissionViewChannel {
					textLockedCount++
					break
				}
			}

			if createInvite && !inviteSet {
				inviteSet = true

				invite, err = common.BotSession.ChannelInviteCreate(c.ID, discordgo.Invite{
					MaxAge:    120,
					MaxUses:   1,
					Temporary: true,
					Unique:    true,
				})
				if err != nil {
					inviteErr = true
					fmt.Printf("INVITE ERROR: %#v", err)
				}
			}
		case discordgo.ChannelTypeGuildVoice:
			voiceCount++
			for _, ow := range c.PermissionOverwrites {
				if ow.Type == discordgo.PermissionOverwriteTypeRole && ow.ID == everyoneID && (ow.Deny&discordgo.PermissionViewChannel) == discordgo.PermissionViewChannel {
					voiceLockedCount++
					break
				}
			}
		}
	}

	if total == 0 {
		channelOutput = "<:red:806413641078407188> No channels"
	} else {
		if categoryCount != 0 {
			channelOutput += fmt.Sprintf("<:list:806420411767586826> %d categories\n", categoryCount)
		}

		if textCount != 0 {
			channelOutput += fmt.Sprintf("<:hashtag:806376515364847626> %d (%d locked)\n", textCount, textLockedCount)
		}

		if voiceCount != 0 {
			channelOutput += fmt.Sprintf("<:mic:806376545936998440> %d (%d locked)\n", voiceCount, voiceLockedCount)
		}

		if newsCount != 0 {
			channelOutput += fmt.Sprintf("<:megaphone:806420565526446100> %d\n", newsCount)
		}

		if storeCount != 0 {
			channelOutput += fmt.Sprintf("<:shop1:806377785823330325>  %d\n", storeCount)
		}
	}

	emojiCount := len(guild.Emojis)
	var animatedEmojiCount int
	var emojiOut string
	for _, e := range guild.Emojis {
		if e.Animated {
			animatedEmojiCount++
		}
	}

	if emojiCount == 0 {
		emojiOut = "<:red:806413641078407188> No custom emojis"
	} else {
		emojiOut = fmt.Sprintf("%d non-animated emojis", emojiCount-animatedEmojiCount)

		if animatedEmojiCount != 0 {
			emojiOut += fmt.Sprintf("\n%d animated emojis", animatedEmojiCount)
		}
	}

	var extras string
	if guild.IconURL() != "" {
		extras = fmt.Sprintf("Guild Region: %s\n[Guild Icon](%s)", strings.Title(guild.Region), guild.IconURL())
	} else {
		extras = fmt.Sprintf("Guild Region: %s", strings.Title(guild.Region))
	}

	verificationLevel := [5]string{"None", "Low", "Medium", "High", "Very High"}
	verificationLevelOut := verificationLevel[int(guild.VerificationLevel)]

	var widgetOut string
	if guild.WidgetEnabled {
		widgetOut = "<:red:806413641078407188>"
	} else {
		widgetOut = "<:green:806413632912490507>"
	}

	roleCount := len(guild.Roles)

	memberCount := len(guild.Members)
	var botCount int
	for _, m := range guild.Members {
		if m != nil && m.User != nil && m.User.Bot {
			botCount++
		}
	}

	var mobileCount, deskptopCount, webCount int
	for _, p := range guild.Presences {
		switch {
		case p.ClientStatus.Desktop != "":
			deskptopCount++
		case p.ClientStatus.Mobile != "":
			mobileCount++
		case p.ClientStatus.Web != "":
			webCount++
		}
	}

	memberOut := fmt.Sprintf("Total: %d\nHumans: %d\nBots: %d", memberCount, memberCount-botCount, botCount)

	var mfaOut string
	if guild.MfaLevel != 0 {
		mfaOut = "<:green:806413632912490507>"
	} else {
		mfaOut = "<:red:806413641078407188>"
	}

	dPremiumCount := guild.PremiumSubscriptionCount
	dPremiumTier := guild.PremiumTier
	explicitContent := [3]string{"<:red:806413641078407188> disabled", "Members Without Roles", "All Members"}
	explicitContentOut := explicitContent[guild.ExplicitContentFilter]

	isPremium, _ := premium.IsGuildPremium(guild.ID)
	var premium string

	if isPremium {
		premium = "<:green:806413632912490507> Premium enabled"
	} else {
		premium = fmt.Sprintf("<:red:806413641078407188> [click here](https://%s/premium)", common.ConfHost.GetString())
	}

	var ownerString string
	owner, err := bot.GetMember(guild.ID, ownerID)
	if err != nil {
		owner = nil
	}

	if owner == nil {
		ownerString = fmt.Sprintf("ID: %d", ownerID)
	} else {
		ownerString = owner.DGoUser().String()
	}

	prefix := common.BotUser.Mention()
	prefix += fmt.Sprintf("\n%s", prfx.GetPrefixIgnoreError(guild.ID))

	created := discordgo.SnowflakeTimestamp(guild.ID)

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       int(rand.Int63n(16777215)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Owner",
				Value:  ownerString,
				Inline: true,
			},
			{
				Name:   "TidBot Premium",
				Value:  premium,
				Inline: true,
			},
			{
				Name:   "TidBot Prefixes",
				Value:  prefix,
				Inline: true,
			},
			{
				Name:   "Features",
				Value:  feat,
				Inline: true,
			},
			{
				Name:   "Custom emojis",
				Value:  emojiOut,
				Inline: true,
			},
			{
				Name:   "WidgetEnabled",
				Value:  widgetOut,
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  memberOut,
				Inline: true,
			},
			{
				Name:   "Roles",
				Value:  fmt.Sprintf("%d roles", roleCount),
				Inline: true,
			},
			{
				Name:   "Channels",
				Value:  channelOutput,
				Inline: true,
			},
			{
				Name:   "Verification Level",
				Value:  verificationLevelOut,
				Inline: true,
			},
			{
				Name:   "Extras",
				Value:  extras,
				Inline: true,
			},
			{
				Name:   "MFA Enabled",
				Value:  mfaOut,
				Inline: true,
			},
			{
				Name:   "Premium Tier",
				Value:  fmt.Sprint(dPremiumTier),
				Inline: true,
			},
			{
				Name:   "Nitro Boosters",
				Value:  fmt.Sprintf("%d boosters", dPremiumCount),
				Inline: true,
			},
			{
				Name:   "Explicit Content Filter",
				Value:  explicitContentOut,
				Inline: true,
			},
			{
				Name:   "Desktop Connections",
				Value:  fmt.Sprint(deskptopCount),
				Inline: true,
			},
			{
				Name:   "Mobile Connections",
				Value:  fmt.Sprint(mobileCount),
				Inline: true,
			},
			{
				Name:   "Web Connections",
				Value:  fmt.Sprint(webCount),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %d, Created", guild.ID),
		},
		Timestamp: created.Format(time.RFC3339),
	}

	if guild.IconURL() != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL(),
		}
	}

	if createInvite {
		if !inviteErr && invite != nil {
			embed.Description += "\n\n**Invite link**: https://discord.gg/" + invite.Code
		} else {
			embed.Description += "\nSomething went wrong while creating the invite link, check console."
		}
	}

	return embed
}
