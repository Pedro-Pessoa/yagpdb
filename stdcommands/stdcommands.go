package stdcommands

import (
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/advice"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/allocstat"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/banserver"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/calc"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/catfact"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/ccreqs"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/covidstats"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/createinvite"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/currentshard"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/currenttime"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/customembed"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/dcallvoice"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/define"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/dogfact"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/editrole"
	edtn "github.com/Pedro-Pessoa/tidbot/stdcommands/edittargetnickname"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/exportcustomcommands"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/findserver"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/getiplocation"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/globalrl"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/guildunavailable"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/info"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/invite"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/leaveserver"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/listroles"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/memberfetcher"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/memstats"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/ping"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/poll"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/roll"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/serverinfo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/setstatus"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/simpleembed"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/sleep"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/stateinfo"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/throw"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/tibiacmds"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/tidstatus"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/toggledbg"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/topcommands"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/topevents"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/topgames"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/topic"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/topservers"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/unbanserver"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/undelete"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/viewperms"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/weather"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/wouldyourather"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/xkcd"
)

var (
	_ bot.BotInitHandler       = (*Plugin)(nil)
	_ commands.CommandProvider = (*Plugin)(nil)
)

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Standard Commands",
		SysName:  "standard_commands",
		Category: common.PluginCategoryCore,
	}
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p,
		// Info
		info.Command,
		invite.Command,

		// Standard
		advice.Command,
		calc.Command,
		catfact.Command,
		currenttime.Command,
		define.Command,
		dogfact.Command,
		topic.Command,
		weather.Command,
		wouldyourather.Command,
		xkcd.Command,
		ping.Command,
		customembed.Command,
		simpleembed.Command,
		editrole.Command,
		listroles.Command,
		memstats.Command,
		poll.Command,
		roll.Command,
		undelete.Command,
		viewperms.Command,
		topgames.Command,
		exportcustomcommands.Command,
		getiplocation.Command,
		covidstats.Command,
		edtn.Command,
		throw.Command,
		serverinfo.Command,

		//Tibia
		//Chars
		tibiacmds.MainCharCommand,
		tibiacmds.DeathsCommand,
		tibiacmds.CheckOnlineCommand,

		//News
		tibiacmds.NewsCommand,
		tibiacmds.NewsTickerCommand,

		//Guilds
		tibiacmds.SpecificGuildCommand,

		//GlobalServerValues
		tibiacmds.TibiaSetWorld,
		tibiacmds.TibiaSetGuild,
		tibiacmds.TibiaGetWorld,
		tibiacmds.TibiaGetGuild,
		tibiacmds.TibiaSetDeathChannel,
		tibiacmds.TibiaSetUpdatesChannel,
		tibiacmds.TibiaCreateNewsFeed,
		tibiacmds.TibiaEnableNewsFeed,
		tibiacmds.TibiaDisableNewsFeed,
		tibiacmds.TibiaChangeNewsFeedChannel,

		//Bot Owner Commands
		tibiacmds.TibiaDelWorld,
		tibiacmds.TibiaDelGuild,
		tibiacmds.TibiaAdmSetWorld,
		tibiacmds.TibiaAdmSetGuild,
		tibiacmds.AdminTrackCommand,
		tibiacmds.AdminTrackHuntedCommand,
		tibiacmds.AdminUntrackHuntedCommand,
		tibiacmds.AdminUntrackGuildCommand,
		tibiacmds.AdminUntrackCommand,
		tibiacmds.AdminDelAllCommand,
		tibiacmds.AdminStartTrackingCommand,
		tibiacmds.AdminStopTrackingCommand,
		tibiacmds.AdminDeleteTracksCommand,
		tibiacmds.AdminStartNewsLoop,
		tibiacmds.AdminStopNewsLoop,
		tibiacmds.AdminDisableNewsFeed,
		tibiacmds.AdminEnableNewsFeed,
		tibiacmds.AdminDebugNewsFeed,

		//Tracking Commands
		tibiacmds.TrackCommand,
		tibiacmds.UnTrackCommand,
		tibiacmds.TrackHuntedCommand,
		tibiacmds.UnTrackHuntedCommand,
		tibiacmds.UnTrackGuildCommand,

		// Maintenance
		stateinfo.Command,
		leaveserver.Command,
		banserver.Command,
		allocstat.Command,
		unbanserver.Command,
		topservers.Command,
		topcommands.Command,
		topevents.Command,
		currentshard.Command,
		memberfetcher.Command,
		guildunavailable.Command,
		tidstatus.Command,
		setstatus.Command,
		createinvite.Command,
		findserver.Command,
		dcallvoice.Command,
		ccreqs.Command,
		sleep.Command,
		toggledbg.Command,
		globalrl.Command,
		serverinfo.AdminCommand,
	)

}

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, ping.HandleMessageCreate, eventsystem.EventMessageCreate)
}

func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{})
}
