package tibia

import (
	"html/template"
	"strconv"

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
	var out string

	for i, w := range TibiaWorlds {
		out += `<option value="` + strconv.Itoa(i) + `">` + w + "</option>\n"
	}

	return template.HTML(out)
}
