package tibia

import (
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/web"
)

var _ web.Plugin = (*Plugin)(nil)

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../tibia/assets/tibia.html", "templates/plugins/tibia.html")
	web.AddSidebarItem(web.SidebarCategoryExtras, &web.SidebarItem{
		Name:   "Tibia",
		NamePT: "Tibia",
		URL:    "tibia",
		Icon:   "fas fa-scroll",
	})

	tibiaMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/tibia"), tibiaMux)
	web.CPMux.Handle(pat.New("/tibia/*"), tibiaMux)

	tibiaMux.Use(web.NotFound())

	tibiaGetHandler := web.RenderHandler(HandleTibiaHtml, "tibia")
	tibiaMux.Handle(pat.Get(""), tibiaGetHandler)
	tibiaMux.Handle(pat.Get("/"), tibiaGetHandler)
}

func HandleTibiaHtml(w http.ResponseWriter, r *http.Request) interface{} {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	tibia := getTibiaData(ag.ID)

	templateData["Tibia"] = tibia

	return templateData
}

func getTibiaData(g int64) map[string]interface{} {
	out := make(map[string]interface{})

	world, _ := GetServerWorld(g, true)
	out["World"] = world

	return out
}
