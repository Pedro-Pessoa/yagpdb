package throw

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/commands"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "Throw",
	Description: "Throwing things is cool.",
	Arguments: []*dcmd.ArgDef{
		{Name: "Target", Type: dcmd.User},
	},

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		target := "a random person nearby"
		if data.Args[0].Value != nil {
			target = data.Args[0].Value.(*discordgo.User).Username
		}

		resp := ""
		rand.Seed(time.Now().UnixNano())
		rng := rand.Intn(100)

		switch {
		case rng < 5:
			resp = fmt.Sprintf("TRIPLE THROW! Threw **%s**, **%s** and **%s** at **%s**", randomThing(), randomThing(), randomThing(), target)
		case rng < 15:
			resp = fmt.Sprintf("DOUBLE THROW! Threw **%s** and **%s** at **%s**", randomThing(), randomThing(), target)
		default:
			resp = fmt.Sprintf("Threw **%s** at **%s**", randomThing(), target)
		}

		return resp, nil
	},
}
