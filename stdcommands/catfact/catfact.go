package catfact

import (
	"math/rand"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryFun,
	Name:        "CatFact",
	Aliases:     []string{"cf", "cat", "catfacts"},
	Description: "Cat Facts",

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return Catfacts[rand.Intn(len(Catfacts))], nil
	},
}
