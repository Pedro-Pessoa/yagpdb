package weather

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lunixbochs/vtclean"

	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

var TempRangeRegex = regexp.MustCompile("(-?[0-9]{1,3})( ?- ?(-?[0-9]{1,3}))? ?°C")

var Command = &commands.TIDCommand{
	CmdCategory:  commands.CategoryFun,
	Name:         "Weather",
	Aliases:      []string{"w"},
	Description:  "Shows the weather somewhere",
	RunInDM:      true,
	RequiredArgs: 1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Where", Type: dcmd.String},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		where := data.Args[0].Str()

		resp, err := http.DefaultClient.Get("http://wttr.in/" + where + "?m")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// remove escape sequences
		unescaped := vtclean.Clean(string(body), false)

		split := strings.Split(string(unescaped), "\n")

		// Show both celcius and fahernheit
		for i, v := range split {
			if !strings.Contains(v, "°C") {
				continue
			}

			var tmpFrom, tmpTo int
			var isRange bool

			submatches := TempRangeRegex.FindStringSubmatch(v)
			if len(submatches) < 2 {
				continue
			}

			tmpFrom, _ = strconv.Atoi(submatches[1])

			if len(submatches) >= 4 && submatches[3] != "" {
				tmpTo, _ = strconv.Atoi(submatches[3])
				isRange = true
			}

			// convert to fahernheit
			tmpFrom = int(float64(tmpFrom)*1.8 + 32)
			tmpTo = int(float64(tmpTo)*1.8 + 32)

			v = strings.TrimRight(v, " ")
			if isRange {
				split[i] = v + " (" + strconv.Itoa(tmpFrom) + "-" + strconv.Itoa(tmpTo) + " °F)"
			} else {
				split[i] = v + " (" + strconv.Itoa(tmpFrom) + " °F)"
			}
		}

		var out strings.Builder
		out.WriteString("```\n")

		for i := 0; i < 7; i++ {
			if i >= len(split) {
				break
			}

			out.WriteString(strings.TrimRight(split[i], " ") + "\n")
		}

		out.WriteString("\n```")

		return out.String(), nil
	},
}
