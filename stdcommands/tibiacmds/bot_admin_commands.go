package tibiacmds

import (
	"fmt"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/premium"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
	"github.com/Pedro-Pessoa/tidbot/tibia"
)

var TibiaDelWorld = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "TibiaDelWorld",
	Aliases:              []string{"tdw"},
	Description:          "Apaga o mundo do servidor.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[0].Value != nil {
			server = data.Args[0].Int64()
		}

		out, err := tibia.DeleteServerWorld(server)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var TibiaDelGuild = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "TibiaDelGuild",
	Aliases:              []string{"tdg"},
	Description:          "Apaga a guild do servidor.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[0].Value != nil {
			server = data.Args[0].Int64()
		}

		out, err := tibia.DeleteServerGuild(server)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var TibiaAdmSetWorld = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "TibiaAdminSetWorld",
	Aliases:              []string{"tasw"},
	Description:          "Reseta o mundo do servidor para o novo.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Mundo", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		out, err := tibia.SetServerWorld(data.Args[0].Str(), server, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var TibiaAdmSetGuild = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "TibiaAdminSetGuild",
	Aliases:              []string{"tasg"},
	Description:          "Reseta a guild do servidor para a nova.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Guild", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		out, err := tibia.SetServerGuild(data.Args[0].Str(), server, true, data.GS.Guild.MemberCount)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminTrackCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "Admintrack",
	Aliases:              []string{"at"},
	Description:          "Faz com que o char especificado seja acompanhado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		isPremium, _ := premium.IsGuildPremium(server)
		out, err := tibia.TrackChar(data.Args[0].Str(), data.GS.ID, data.GS.Guild.MemberCount, isPremium, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminTrackHuntedCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "Admintrackhunted",
	Aliases:              []string{"ath"},
	Description:          "Faz com que o char especificado seja acompanhado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		isPremium, _ := premium.IsGuildPremium(server)
		out, err := tibia.TrackChar(data.Args[0].Str(), data.GS.ID, data.GS.Guild.MemberCount, isPremium, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminUntrackCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminUntrack",
	Aliases:              []string{"au"},
	Description:          "Faz com que o char especificado deixe de ser acompanhado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		out, err := tibia.UnTrackChar(data.Args[0].Str(), server, false, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminUntrackHuntedCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminUntrackHunted",
	Aliases:              []string{"auh"},
	Description:          "Faz com que o char especificado deixe de ser acompanhado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		out, err := tibia.UnTrackChar(data.Args[0].Str(), server, true, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminUntrackGuildCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminUntrackGuild",
	Aliases:              []string{"aug"},
	Description:          "Faz com que o char especificado deixe de ser acompanhado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[1].Int64()
		}

		out, err := tibia.UnTrackChar(data.Args[0].Str(), server, false, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminDelAllCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminDelAll",
	Aliases:              []string{"ada"},
	Description:          "Deleta TODAS as databases de tibia.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.DeleteAll()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminStartTrackingCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminStartTracking",
	Aliases:              []string{"ast", "tracking", "starttrack"},
	Description:          "Inicia o tracking.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.StartLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminStopTrackingCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminStopTracking",
	Aliases:              []string{"astopt", "stoptrack"},
	Description:          "Para o tracking.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.StopLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminDeleteTracksCommand = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminDeleteTracks",
	Aliases:              []string{"adt"},
	Description:          "Deleta o track especificado do server especificado.",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "Track a deletar", Type: dcmd.String},
		{Name: "ID do servidor.", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[1].Value != nil {
			server = data.Args[0].Int64()
		}
		switch data.Args[0].Str() {
		case "all":
			out, err := tibia.DeleteTracks(server, false, false, true)
			if err != nil {
				return fmt.Sprintln(err), err
			}

			return out, nil
		case "hunted", "hunteds":
			out, err := tibia.DeleteTracks(server, true, false, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}

			return out, nil
		case "guild", "guilds":
			out, err := tibia.DeleteTracks(server, false, true, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}

			return out, nil
		case "track", "tracks":
			out, err := tibia.DeleteTracks(server, false, false, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}

			return out, nil
		default:
			return "Track inválido.", nil
		}
	}),
}

var AdminStartNewsLoop = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "StartNewsLoop",
	Aliases:              []string{"snl"},
	Description:          "Inicia o loop de news",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.StartNewsLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminStopNewsLoop = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "StopNewsLoop",
	Aliases:              []string{"stnl"},
	Description:          "Para o loop de news",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.StopNewsLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminDisableNewsFeed = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminDisableNewsFeed",
	Aliases:              []string{"adnf"},
	Description:          "Para o loop de news no servidor (pode targetar um servidor opcionalmente)",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "id", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[0].Value != nil {
			server = data.Args[0].Int64()
		}

		out, err := tibia.DisableNewsFeed(server)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminEnableNewsFeed = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminEnableNewsFeed",
	Aliases:              []string{"aenf"},
	Description:          "Inicia o loop de news no servidor (pode targetar um servidor opcionalmente)",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	Arguments: []*dcmd.ArgDef{
		{Name: "id", Type: dcmd.Int},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		server := data.GS.ID
		if data.Args[0].Value != nil {
			server = data.Args[0].Int64()
		}

		out, err := tibia.EnableNewsFeed(server)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}

var AdminDebugNewsFeed = &commands.TIDCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminDebugNewsFeed",
	Aliases:              []string{"adnf", "adf"},
	Description:          "Debug News Feed Cmd",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.DebugNews(data.CS.ID)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	}),
}
