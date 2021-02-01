package toggledbg

import (
	"github.com/sirupsen/logrus"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/stdcommands/util"
)

var Command = &commands.TIDCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "toggledbg",
	Description:          "Toggles Debug Logging",
	HideFromHelp:         true,
	IsModCmd:             true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			common.SetLoggingLevel(logrus.InfoLevel)
			return "Disabled debug logging", nil
		}

		common.SetLoggingLevel(logrus.DebugLevel)
		return "Enabled debug logging", nil
	}),
}
