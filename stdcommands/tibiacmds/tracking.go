package tibiacmds

import (
	"fmt"

	"emperror.dev/errors"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/premium"
	"github.com/jonas747/yagpdb/stdcommands/util"
	"github.com/jonas747/yagpdb/tibia"
)

var TrackCommand = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "Track",
	Description:  "Faz com que o char especificado seja acompanhado.",
	RequiredArgs: 1,
	Cooldown:     30,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		isPremium, _ := premium.IsGuildPremium(data.GS.ID)
		out, err := tibia.TrackChar(data.Args[0].Str(), data.GS.ID, data.GS.Guild.MemberCount, isPremium, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var TrackHuntedCommand = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "TrackHunted",
	Description:  "Faz com que o char especificado seja acompanhado.",
	RequiredArgs: 1,
	Cooldown:     30,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if util.IsExecedByCC(data) {
			return "", errors.New("Esse comando não pode ser executado através de um Custom Command.")
		}

		isPremium, _ := premium.IsGuildPremium(data.GS.ID)
		out, err := tibia.TrackChar(data.Args[0].Str(), data.GS.ID, data.GS.Guild.MemberCount, isPremium, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var UnTrackCommand = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "Untrack",
	Description:  "Faz com que o char especificado deixe de ser acompanhado.",
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.UnTrackChar(data.Args[0].Str(), data.GS.ID, false, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var UnTrackHuntedCommand = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "Untrackhunted",
	Description:  "Faz com que o char especificado deixe de ser acompanhado.",
	Aliases:      []string{"uth"},
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.UnTrackChar(data.Args[0].Str(), data.GS.ID, true, false)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}

var UnTrackGuildCommand = &commands.YAGCommand{
	CmdCategory:  commands.CategoryTibia,
	Name:         "Untrackguild",
	Description:  "Faz com que o char especificado deixe de ser acompanhado.",
	Aliases:      []string{"utg"},
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Nome do Char", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		out, err := tibia.UnTrackChar(data.Args[0].Str(), data.GS.ID, false, true)
		if err != nil {
			return fmt.Sprintln(err), err
		}

		return out, nil
	},
}
