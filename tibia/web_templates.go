package tibia

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/Pedro-Pessoa/tidbot/web"
)

var tibiaTemplates map[string]interface{}

func init() {
	tibiaTemplates = insertTemplates()
	web.RegisterSetupFuncWeb(tibiaTemplates)
}

func insertTemplates() map[string]interface{} {
	out := map[string]interface{}{
		"tibiaWorldDropDown": tmplTibiaWorldDropDown,
	}

	return out
}

func tmplTibiaWorldDropDown() template.HTML {
	var out strings.Builder

	for i, w := range TibiaWorlds {
		out.WriteString(`<option value="` + strconv.Itoa(i) + `">` + w + "</option>\n")
	}

	return template.HTML(out.String())
}
