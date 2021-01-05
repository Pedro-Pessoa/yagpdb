package invite

import (
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Invite",
	Description: "Responds with bot invite link",
	RunInDM:     true,

	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		return "Por favor adicione o bot através do site\nhttps://" + common.ConfHost.GetString(), nil
	},
}
