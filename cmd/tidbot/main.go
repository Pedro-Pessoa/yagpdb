package main

import (
	"github.com/Pedro-Pessoa/tidbot/analytics"
	"github.com/Pedro-Pessoa/tidbot/common/featureflags"
	"github.com/Pedro-Pessoa/tidbot/common/prom"
	"github.com/Pedro-Pessoa/tidbot/common/run"
	"github.com/Pedro-Pessoa/tidbot/web/discorddata"

	// Core yagpdb packages

	"github.com/Pedro-Pessoa/tidbot/admin"
	"github.com/Pedro-Pessoa/tidbot/bot/paginatedmessages"
	"github.com/Pedro-Pessoa/tidbot/common/internalapi"
	"github.com/Pedro-Pessoa/tidbot/common/scheduledevents2"

	// Plugin imports
	"github.com/Pedro-Pessoa/tidbot/automod"
	"github.com/Pedro-Pessoa/tidbot/automod_legacy"
	"github.com/Pedro-Pessoa/tidbot/autorole"
	"github.com/Pedro-Pessoa/tidbot/aylien"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/customcommands"
	"github.com/Pedro-Pessoa/tidbot/discordlogger"
	"github.com/Pedro-Pessoa/tidbot/logs"
	"github.com/Pedro-Pessoa/tidbot/moderation"
	"github.com/Pedro-Pessoa/tidbot/notifications"
	"github.com/Pedro-Pessoa/tidbot/premium"
	"github.com/Pedro-Pessoa/tidbot/premium/patreonpremiumsource"
	"github.com/Pedro-Pessoa/tidbot/reddit"
	"github.com/Pedro-Pessoa/tidbot/reminders"
	"github.com/Pedro-Pessoa/tidbot/reputation"
	"github.com/Pedro-Pessoa/tidbot/rolecommands"
	"github.com/Pedro-Pessoa/tidbot/safebrowsing"
	"github.com/Pedro-Pessoa/tidbot/serverstats"
	"github.com/Pedro-Pessoa/tidbot/stdcommands"
	"github.com/Pedro-Pessoa/tidbot/streaming"
	"github.com/Pedro-Pessoa/tidbot/tibia"
	"github.com/Pedro-Pessoa/tidbot/tickets"
	"github.com/Pedro-Pessoa/tidbot/twitter"
	"github.com/Pedro-Pessoa/tidbot/verification"
	"github.com/Pedro-Pessoa/tidbot/youtube"
	// External plugins
)

func main() {

	run.Init()

	//BotSession.LogLevel = discordgo.LogInformational
	paginatedmessages.RegisterPlugin()
	discorddata.RegisterPlugin()

	// Setup plugins
	analytics.RegisterPlugin()
	safebrowsing.RegisterPlugin()
	discordlogger.Register()
	commands.RegisterPlugin()
	stdcommands.RegisterPlugin()
	serverstats.RegisterPlugin()
	notifications.RegisterPlugin()
	customcommands.RegisterPlugin()
	reddit.RegisterPlugin()
	moderation.RegisterPlugin()
	tibia.RegisterPlugin()
	reputation.RegisterPlugin()
	aylien.RegisterPlugin()
	streaming.RegisterPlugin()
	automod_legacy.RegisterPlugin()
	automod.RegisterPlugin()
	logs.RegisterPlugin()
	autorole.RegisterPlugin()
	reminders.RegisterPlugin()
	youtube.RegisterPlugin()
	rolecommands.RegisterPlugin()
	tickets.RegisterPlugin()
	verification.RegisterPlugin()
	premium.RegisterPlugin()
	patreonpremiumsource.RegisterPlugin()
	scheduledevents2.RegisterPlugin()
	twitter.RegisterPlugin()
	admin.RegisterPlugin()
	internalapi.RegisterPlugin()
	prom.RegisterPlugin()
	featureflags.RegisterPlugin()

	run.Run()
}
