package tibia

import (
	"fmt"
	"html/template"

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
	out := ""

	for i, w := range tibiaWorlds {
		out += `<option value="` + fmt.Sprint(i) + `">` + w + "</option>\n"
	}

	return template.HTML(out)
}

var tibiaWorlds = []string{
	"Antica",
	"Ferobra",
	"Serdebra",
}
