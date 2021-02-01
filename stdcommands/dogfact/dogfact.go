package dogfact

import (
	"math/rand"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "DogFact",
	Aliases:     []string{"dog", "dogfacts"},
	Description: "Dog Facts",

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return dogfacts[rand.Intn(len(dogfacts))], nil
	},
}
