package topcommands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	Cooldown:     2,
	CmdCategory:  commands.CategoryDebug,
	Name:         "topcommands",
	Description:  "Shows command usage stats",
	HideFromHelp: true,
	Arguments: []*dcmd.ArgDef{
		{Name: "hours", Type: dcmd.Int, Default: 1},
	},
	RunFunc: cmdFuncTopCommands,
}

func cmdFuncTopCommands(data *dcmd.Data) (interface{}, error) {
	hours := data.Args[0].Int()
	within := time.Now().Add(time.Duration(-hours) * time.Hour)

	var results []*TopCommandsResult
	err := common.GORM.Table(common.LoggedExecutedCommand{}.TableName()).Select("command, COUNT(id)").Where("created_at > ?", within).Group("command").Order("count(id) desc").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var out strings.Builder
	out.WriteString("```\nCommand stats from now to " + strconv.Itoa(hours) + "hour(s) ago\n#    Total -  Command\n")
	var total int

	for k, result := range results {
		out.WriteString(fmt.Sprintf("#%02s: %5s - %s\n", strconv.Itoa(k+1), strconv.Itoa(result.Count), result.Command))
		total += result.Count
	}

	cpm := float64(total) / float64(hours) / 60

	out.WriteString(fmt.Sprintf("\nTotal: %s, Commands per minute: %.1f\n```", strconv.Itoa(total), cpm))

	return out.String(), nil
}

type TopCommandsResult struct {
	Command string
	Count   int
}
