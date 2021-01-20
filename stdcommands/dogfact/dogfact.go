package dogfact

import (
	"math/rand"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "DogFact",
	Aliases:     []string{"dog", "dogfacts"},
	Description: "Dog Facts",

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return dogfacts[rand.Intn(len(dogfacts))], nil
	},
}
