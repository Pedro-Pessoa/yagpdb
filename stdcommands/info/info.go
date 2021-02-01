package info

import (
	"fmt"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Info",
	Description: "Responds with bot information",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		info := fmt.Sprintf(`TID <https://%s/manage>`, common.ConfHost.GetString())

		return info, nil
	},
}
