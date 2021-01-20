package roll

import (
	"fmt"
	"math/rand"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/dice"
	"github.com/jonas747/yagpdb/commands"
)

var Command = &commands.YAGCommand{
	CmdCategory:     commands.CategoryFun,
	Name:            "Roll",
	Description:     "Roll dices, specify nothing for 6 sides, specify a number for max sides, or rpg dice syntax.",
	LongDescription: "Example: `-roll 2d6`",
	Arguments: []*dcmd.ArgDef{
		{Name: "RPG Dice", Type: dcmd.String},
		{Name: "Sides", Default: 0, Type: dcmd.Int},
	},
	ArgumentCombos: [][]int{{1}, {0}, {}},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		if data.Args[0].Value != nil {
			// Special dice syntax if string
			r, _, err := dice.Roll(data.Args[0].Str())
			if err != nil {
				return err.Error(), nil
			}

			output := r.String()
			if len(output) > 100 {
				output = output[:100] + "..."
			}
			return output, nil
		}

		// normal, n sides dice rolling
		sides := data.Args[1].Int()
		if sides < 1 {
			sides = 6
		}

		return fmt.Sprintf(":game_die: %d (1 - %d)", rand.Intn(sides)+1, sides), nil
	},
}
