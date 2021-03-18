// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains variables for all known Discord end points.  All functions
// throughout the Discordgo package use these variables for all connections
// to Discord.  These are all exported and you may modify them if needed.

package discordgo

import "strconv"

// APIVersion is the Discord API version used for the REST and Websocket API.
var APIVersion = "8"

// Known Discord API Endpoints.
var (
	EndpointStatus     = "https://status.discord.com/api/v2/"
	EndpointSm         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive   = EndpointSm + "active.json"
	EndpointSmUpcoming = EndpointSm + "upcoming.json"

	EndpointDiscord    = "https://discord.com/"
	EndpointAPI        = EndpointDiscord + "api/v" + APIVersion + "/"
	EndpointGuilds     = EndpointAPI + "guilds/"
	EndpointChannels   = EndpointAPI + "channels/"
	EndpointUsers      = EndpointAPI + "users/"
	EndpointGateway    = EndpointAPI + "gateway"
	EndpointGatewayBot = EndpointGateway + "/bot"
	EndpointWebhooks   = EndpointAPI + "webhooks/"

	EndpointCDN             = "https://cdn.discordapp.com/"
	EndpointCDNAttachments  = EndpointCDN + "attachments/"
	EndpointCDNAvatars      = EndpointCDN + "avatars/"
	EndpointCDNIcons        = EndpointCDN + "icons/"
	EndpointCDNSplashes     = EndpointCDN + "splashes/"
	EndpointCDNChannelIcons = EndpointCDN + "channel-icons/"
	EndpointCDNBanners      = EndpointCDN + "banners/"

	EndpointAuth           = EndpointAPI + "auth/"
	EndpointLogin          = EndpointAuth + "login"
	EndpointLogout         = EndpointAuth + "logout"
	EndpointVerify         = EndpointAuth + "verify"
	EndpointVerifyResend   = EndpointAuth + "verify/resend"
	EndpointForgotPassword = EndpointAuth + "forgot"
	EndpointResetPassword  = EndpointAuth + "reset"
	EndpointRegister       = EndpointAuth + "register"

	EndpointVoice        = EndpointAPI + "/voice/"
	EndpointVoiceRegions = EndpointVoice + "regions"
	EndpointVoiceIce     = EndpointVoice + "ice"

	EndpointTutorial           = EndpointAPI + "tutorial/"
	EndpointTutorialIndicators = EndpointTutorial + "indicators"

	EndpointTrack        = EndpointAPI + "track"
	EndpointSso          = EndpointAPI + "sso"
	EndpointReport       = EndpointAPI + "report"
	EndpointIntegrations = EndpointAPI + "integrations"

	EndpointUser               = func(uID string) string { return EndpointUsers + uID }
	EndpointUserAvatar         = func(uID int64, aID string) string { return EndpointCDNAvatars + StrID(uID) + "/" + aID + ".png" }
	EndpointUserAvatarAnimated = func(uID int64, aID string) string { return EndpointCDNAvatars + StrID(uID) + "/" + aID + ".gif" }
	EndpointUserSettings       = func(uID string) string { return EndpointUsers + uID + "/settings" }
	EndpointUserGuilds         = func(uID string) string { return EndpointUsers + uID + "/guilds" }
	EndpointUserGuild          = func(uID string, gID int64) string { return EndpointUsers + uID + "/guilds/" + StrID(gID) }
	EndpointUserGuildSettings  = func(uID string, gID int64) string { return EndpointUsers + uID + "/guilds/" + StrID(gID) + "/settings" }
	EndpointUserChannels       = func(uID string) string { return EndpointUsers + uID + "/channels" }
	EndpointUserDevices        = func(uID string) string { return EndpointUsers + uID + "/devices" }
	EndpointUserConnections    = func(uID string) string { return EndpointUsers + uID + "/connections" }
	EndpointUserNotes          = func(uID int64) string { return EndpointUsers + "@me/notes/" + StrID(uID) }

	EndpointGuild           = func(gID int64) string { return EndpointGuilds + StrID(gID) }
	EndpointGuildPreview    = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/preview" }
	EndpointGuildChannels   = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/channels" }
	EndpointGuildMembers    = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/members" }
	EndpointGuildMember     = func(gID int64, uID int64) string { return EndpointGuilds + StrID(gID) + "/members/" + StrID(uID) }
	EndpointGuildMemberMe   = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/members/@me" }
	EndpointGuildMemberRole = func(gID, uID, rID int64) string {
		return EndpointGuilds + StrID(gID) + "/members/" + StrID(uID) + "/roles/" + StrID(rID)
	}
	EndpointGuildBans            = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/bans" }
	EndpointGuildBan             = func(gID, uID int64) string { return EndpointGuilds + StrID(gID) + "/bans/" + StrID(uID) }
	EndpointGuildIntegrations    = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/integrations" }
	EndpointGuildIntegration     = func(gID, iID int64) string { return EndpointGuilds + StrID(gID) + "/integrations/" + StrID(iID) }
	EndpointGuildIntegrationSync = func(gID, iID int64) string {
		return EndpointGuilds + StrID(gID) + "/integrations/" + StrID(iID) + "/sync"
	}
	EndpointGuildRoles        = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/roles" }
	EndpointGuildRole         = func(gID, rID int64) string { return EndpointGuilds + StrID(gID) + "/roles/" + StrID(rID) }
	EndpointGuildInvites      = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/invites" }
	EndpointGuildWidget       = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/widget" }
	EndpointGuildEmbed        = EndpointGuildWidget
	EndpointGuildPrune        = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/prune" }
	EndpointGuildIcon         = func(gID int64, hash string) string { return EndpointCDNIcons + StrID(gID) + "/" + hash + ".png" }
	EndpointGuildIconAnimated = func(gID int64, hash string) string { return EndpointCDNIcons + StrID(gID) + "/" + hash + ".gif" }
	EndpointGuildSplash       = func(gID int64, hash string) string { return EndpointCDNSplashes + StrID(gID) + "/" + hash + ".png" }
	EndpointGuildWebhooks     = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/webhooks" }
	EndpointGuildAuditLogs    = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/audit-logs" }
	EndpointGuildEmojis       = func(gID int64) string { return EndpointGuilds + StrID(gID) + "/emojis" }
	EndpointGuildEmoji        = func(gID, eID int64) string { return EndpointGuilds + StrID(gID) + "/emojis/" + StrID(eID) }
	EndpointGuildBanner       = func(gID int64, hash string) string { return EndpointCDNBanners + StrID(gID) + "/" + hash + ".png" }

	EndpointChannel                   = func(cID int64) string { return EndpointChannels + StrID(cID) }
	EndpointChannelPermissions        = func(cID int64) string { return EndpointChannels + StrID(cID) + "/permissions" }
	EndpointChannelPermission         = func(cID, tID int64) string { return EndpointChannels + StrID(cID) + "/permissions/" + StrID(tID) }
	EndpointChannelInvites            = func(cID int64) string { return EndpointChannels + StrID(cID) + "/invites" }
	EndpointChannelTyping             = func(cID int64) string { return EndpointChannels + StrID(cID) + "/typing" }
	EndpointChannelMessages           = func(cID int64) string { return EndpointChannels + StrID(cID) + "/messages" }
	EndpointChannelMessage            = func(cID, mID int64) string { return EndpointChannels + StrID(cID) + "/messages/" + StrID(mID) }
	EndpointChannelMessageAck         = func(cID, mID int64) string { return EndpointChannels + StrID(cID) + "/messages/" + StrID(mID) + "/ack" }
	EndpointChannelMessagesBulkDelete = func(cID int64) string { return EndpointChannel(cID) + "/messages/bulk-delete" }
	EndpointChannelMessagesPins       = func(cID int64) string { return EndpointChannel(cID) + "/pins" }
	EndpointChannelMessagePin         = func(cID, mID int64) string { return EndpointChannel(cID) + "/pins/" + StrID(mID) }
	EndpointChannelMessageCrosspost   = func(cID, mID int64) string { return EndpointChannel(cID) + "/messages/" + StrID(mID) + "/crosspost" }
	EndpointChannelFollow             = func(cID int64) string { return EndpointChannel(cID) + "/followers" }

	EndpointGroupIcon = func(cID int64, hash string) string { return EndpointCDNChannelIcons + StrID(cID) + "/" + hash + ".png" }

	EndpointChannelWebhooks = func(cID int64) string { return EndpointChannel(cID) + "/webhooks" }
	EndpointWebhook         = func(wID int64) string { return EndpointWebhooks + StrID(wID) }
	EndpointWebhookToken    = func(wID int64, token string) string { return EndpointWebhooks + StrID(wID) + "/" + token }
	EndpointWebhookMessage  = func(wID int64, token, messageID string) string {
		return EndpointWebhookToken(wID, token) + "/messages/" + messageID
	}

	EndpointDefaultUserAvatar = func(uDiscriminator string) string {
		uDiscriminatorInt, _ := strconv.Atoi(uDiscriminator)
		return EndpointCDN + "embed/avatars/" + strconv.Itoa(uDiscriminatorInt%5) + ".png"
	}

	EndpointMessageReactionsAll = func(cID, mID int64) string {
		return EndpointChannelMessage(cID, mID) + "/reactions"
	}
	EndpointMessageReactions = func(cID, mID int64, emoji EmojiName) string {
		return EndpointChannelMessage(cID, mID) + "/reactions/" + emoji.String()
	}
	EndpointMessageReaction = func(cID, mID int64, emoji EmojiName, uID string) string {
		return EndpointMessageReactions(cID, mID, emoji) + "/" + uID
	}

	EndpointApplicationGlobalCommands = func(aID int64) string {
		return EndpointApplication(aID) + "/commands"
	}
	EndpointApplicationGlobalCommand = func(aID, cID int64) string {
		return EndpointApplicationGlobalCommands(aID) + "/" + StrID(cID)
	}

	EndpointApplicationGuildCommands = func(aID, gID int64) string {
		return EndpointApplication(aID) + "/guilds/" + StrID(gID) + "/commands"
	}
	EndpointApplicationGuildCommand = func(aID, gID, cID int64) string {
		return EndpointApplicationGuildCommands(aID, gID) + "/" + StrID(cID)
	}
	EndpointInteraction = func(aID int64, iToken string) string {
		return EndpointAPI + "interactions/" + StrID(aID) + "/" + iToken
	}
	EndpointInteractionResponse = func(iID int64, iToken string) string {
		return EndpointInteraction(iID, iToken) + "/callback"
	}
	EndpointInteractionResponseActions = func(aID int64, iToken string) string {
		return EndpointWebhookMessage(aID, iToken, "@original")
	}
	EndpointFollowupMessage = func(aID int64, iToken string) string {
		return EndpointWebhookToken(aID, iToken)
	}
	EndpointFollowupMessageActions = func(aID int64, iToken, mID string) string {
		return EndpointWebhookMessage(aID, iToken, mID)
	}

	EndpointRelationships       = func() string { return EndpointUsers + "@me" + "/relationships" }
	EndpointRelationship        = func(uID int64) string { return EndpointRelationships() + "/" + StrID(uID) }
	EndpointRelationshipsMutual = func(uID int64) string { return EndpointUsers + StrID(uID) + "/relationships" }

	EndpointGuildCreate = EndpointAPI + "guilds"

	EndpointInvite = func(iID string) string { return EndpointAPI + "invites/" + iID }

	EndpointIntegrationsJoin = func(iID string) string { return EndpointAPI + "integrations/" + iID + "/join" }

	EndpointEmoji         = func(eID int64) string { return EndpointCDN + "emojis/" + StrID(eID) + ".png" }
	EndpointEmojiAnimated = func(eID int64) string { return EndpointCDN + "emojis/" + StrID(eID) + ".gif" }

	EndpointApplications = EndpointAPI + "applications"
	EndpointApplication  = func(aID int64) string { return EndpointApplications + "/" + StrID(aID) }

	EndpointOAuth2                  = EndpointAPI + "oauth2/"
	EndpointOAuth2Applications      = EndpointOauth2 + "applications"
	EndpointOAuth2Application       = func(aID int64) string { return EndpointApplications + "/" + StrID(aID) }
	EndpointOAuth2ApplicationsBot   = func(aID int64) string { return EndpointApplications + "/" + StrID(aID) + "/bot" }
	EndpointOAuth2ApplicationAssets = func(aID int64) string { return EndpointApplications + "/" + StrID(aID) + "/assets" }

	// Deprecated
	EndpointOauth2                  = EndpointOAuth2
	EndpointOauth2Applications      = EndpointOAuth2Applications
	EndpointOauth2Application       = EndpointOAuth2Application
	EndpointOauth2ApplicationsBot   = EndpointOAuth2ApplicationsBot
	EndpointOauth2ApplicationAssets = EndpointOAuth2ApplicationAssets
)
