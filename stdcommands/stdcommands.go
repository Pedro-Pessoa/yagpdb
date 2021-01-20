package stdcommands

import (
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/stdcommands/advice"
	"github.com/jonas747/yagpdb/stdcommands/allocstat"
	"github.com/jonas747/yagpdb/stdcommands/banserver"
	"github.com/jonas747/yagpdb/stdcommands/calc"
	"github.com/jonas747/yagpdb/stdcommands/catfact"
	"github.com/jonas747/yagpdb/stdcommands/ccreqs"
	"github.com/jonas747/yagpdb/stdcommands/covidstats"
	"github.com/jonas747/yagpdb/stdcommands/createinvite"
	"github.com/jonas747/yagpdb/stdcommands/currentshard"
	"github.com/jonas747/yagpdb/stdcommands/currenttime"
	"github.com/jonas747/yagpdb/stdcommands/customembed"
	"github.com/jonas747/yagpdb/stdcommands/dcallvoice"
	"github.com/jonas747/yagpdb/stdcommands/define"
	"github.com/jonas747/yagpdb/stdcommands/dogfact"
	"github.com/jonas747/yagpdb/stdcommands/editrole"
	edtn "github.com/jonas747/yagpdb/stdcommands/edittargetnickname"
	"github.com/jonas747/yagpdb/stdcommands/exportcustomcommands"
	"github.com/jonas747/yagpdb/stdcommands/findserver"
	"github.com/jonas747/yagpdb/stdcommands/getiplocation"
	"github.com/jonas747/yagpdb/stdcommands/globalrl"
	"github.com/jonas747/yagpdb/stdcommands/guildunavailable"
	"github.com/jonas747/yagpdb/stdcommands/info"
	"github.com/jonas747/yagpdb/stdcommands/invite"
	"github.com/jonas747/yagpdb/stdcommands/leaveserver"
	"github.com/jonas747/yagpdb/stdcommands/listroles"
	"github.com/jonas747/yagpdb/stdcommands/memberfetcher"
	"github.com/jonas747/yagpdb/stdcommands/memstats"
	"github.com/jonas747/yagpdb/stdcommands/ping"
	"github.com/jonas747/yagpdb/stdcommands/poll"
	"github.com/jonas747/yagpdb/stdcommands/roll"
	"github.com/jonas747/yagpdb/stdcommands/setstatus"
	"github.com/jonas747/yagpdb/stdcommands/simpleembed"
	"github.com/jonas747/yagpdb/stdcommands/sleep"
	"github.com/jonas747/yagpdb/stdcommands/stateinfo"
	"github.com/jonas747/yagpdb/stdcommands/throw"
	"github.com/jonas747/yagpdb/stdcommands/tibiacmds"
	"github.com/jonas747/yagpdb/stdcommands/toggledbg"
	"github.com/jonas747/yagpdb/stdcommands/topcommands"
	"github.com/jonas747/yagpdb/stdcommands/topevents"
	"github.com/jonas747/yagpdb/stdcommands/topgames"
	"github.com/jonas747/yagpdb/stdcommands/topic"
	"github.com/jonas747/yagpdb/stdcommands/topservers"
	"github.com/jonas747/yagpdb/stdcommands/unbanserver"
	"github.com/jonas747/yagpdb/stdcommands/undelete"
	"github.com/jonas747/yagpdb/stdcommands/viewperms"
	"github.com/jonas747/yagpdb/stdcommands/weather"
	"github.com/jonas747/yagpdb/stdcommands/wouldyourather"
	"github.com/jonas747/yagpdb/stdcommands/xkcd"
	"github.com/jonas747/yagpdb/stdcommands/yagstatus"
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
		yagstatus.Command,
		setstatus.Command,
		createinvite.Command,
		findserver.Command,
		dcallvoice.Command,
		ccreqs.Command,
		sleep.Command,
		toggledbg.Command,
		globalrl.Command,
	)

}

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, ping.HandleMessageCreate, eventsystem.EventMessageCreate)
}

func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{})
}
