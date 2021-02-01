package invite

import (
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Invite",
	Description: "Responds with bot invite link",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return "Por favor adicione o bot através do site\nhttps://" + common.ConfHost.GetString(), nil
	},
}
