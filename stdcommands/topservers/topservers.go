package topservers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/Pedro-Pessoa/tidbot/bot/models"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var Command = &commands.TIDCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryFun,
	Name:        "TopServers",
	Description: "Responds with the top 20 servers I'm on",
	Arguments: []*dcmd.ArgDef{
		{Name: "Skip", Help: "Entries to skip", Type: dcmd.Int, Default: 0},
	},
	ArgSwitches: []*dcmd.ArgDef{
		{Switch: "id", Name: "serverID", Type: dcmd.Int},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		skip := data.Args[0].Int()

		if data.Switches["id"].Value != nil {
			type serverIDQuery struct {
				MemberCount int64
				Name        string
				Place       int64
			}

			var position serverIDQuery
			serverID := data.Switch("id").Int64()

			const q = `SELECT member_count, name, row_number FROM (SELECT id, member_count, name, left_at, row_number() OVER (ORDER BY member_count DESC) FROM joined_guilds WHERE left_at IS NULL) AS total WHERE id=$1 AND left_at IS NULL;`
			err := common.PQ.QueryRow(q, serverID).Scan(&position.MemberCount, &position.Name, &position.Place)
			return fmt.Sprintf("```Server with ID %d is placed:\n#%-2d: %-25s (%d members)\n```", serverID, position.Place, position.Name, position.MemberCount), err
		}

		results, err := models.JoinedGuilds(qm.Where("left_at is null"), qm.OrderBy("member_count desc"), qm.Limit(20), qm.Offset(skip)).AllG(data.Context())
		if err != nil {
			return nil, err
		}

		var out strings.Builder
		out.WriteString("```")

		for k, v := range results {
			out.WriteString(fmt.Sprintf("\n#%2s: %25s (%s members)", strconv.Itoa(k+skip+1), v.Name, strconv.FormatInt(v.MemberCount, 10)))
		}

		return "Top servers the bot is on:\n" + out.String() + "\n```", nil
	},
}
