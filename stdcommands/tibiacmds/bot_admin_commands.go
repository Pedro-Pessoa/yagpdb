package tibiacmds

import (
	"fmt"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/premium"
	"github.com/jonas747/yagpdb/stdcommands/util"
	"github.com/jonas747/yagpdb/tibia"
)

var TibiaDelWorld = &commands.YAGCommand{
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

var TibiaDelGuild = &commands.YAGCommand{
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

var TibiaAdmSetWorld = &commands.YAGCommand{
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

var TibiaAdmSetGuild = &commands.YAGCommand{
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

var AdminTrackCommand = &commands.YAGCommand{
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

var AdminTrackHuntedCommand = &commands.YAGCommand{
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
		a, err := tibia.TrackChar(data.Args[0].Str(), data.GS.ID, data.GS.Guild.MemberCount, isPremium, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminUntrackCommand = &commands.YAGCommand{
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
		a, err := tibia.UnTrackChar(data.Args[0].Str(), server, false, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminUntrackHuntedCommand = &commands.YAGCommand{
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
		a, err := tibia.UnTrackChar(data.Args[0].Str(), server, true, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminUntrackGuildCommand = &commands.YAGCommand{
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
		a, err := tibia.UnTrackChar(data.Args[0].Str(), server, false, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminDellAllCommand = &commands.YAGCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminDellAll",
	Aliases:              []string{"ada"},
	Description:          "Deleta TODAS as databases de tibia.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		a, err := tibia.DeleteAll()
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminStartTrackingCommand = &commands.YAGCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminStartTracking",
	Aliases:              []string{"ast", "tracking", "starttrack"},
	Description:          "Inicia o tracking.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		a, err := tibia.StartLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminStopTrackingCommand = &commands.YAGCommand{
	CmdCategory:          commands.CategoryTibia,
	Name:                 "AdminStopTracking",
	Aliases:              []string{"astopt", "stoptrack"},
	Description:          "Para o tracking.",
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		a, err := tibia.StopLoop()
		if err != nil {
			return fmt.Sprintln(err), err
		}
		return a, nil
	}),
}

var AdminDeleteTracksCommand = &commands.YAGCommand{
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
			a, err := tibia.DeleteTracks(server, false, false, true)
			if err != nil {
				return fmt.Sprintln(err), err
			}
			return a, nil
		case "hunted", "hunteds":
			a, err := tibia.DeleteTracks(server, true, false, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}
			return a, nil
		case "guild", "guilds":
			a, err := tibia.DeleteTracks(server, false, true, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}
			return a, nil
		case "track", "tracks":
			a, err := tibia.DeleteTracks(server, false, false, false)
			if err != nil {
				return fmt.Sprintln(err), err
			}
			return a, nil
		default:
			return "Track inválido.", nil
		}
	}),
}
