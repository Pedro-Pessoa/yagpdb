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

	tibiaGetHandler := web.ControllerHandler(publicHandler(HandleTibiaHtml, false), "tibia")
	tibiaMux.Handle(pat.Get(""), tibiaGetHandler)
	tibiaMux.Handle(pat.Get("/"), tibiaGetHandler)
}

type publicHandlerFunc func(w http.ResponseWriter, r *http.Request, publicAccess bool) (web.TemplateData, error)

func publicHandler(inner publicHandlerFunc, public bool) web.ControllerHandlerFunc {
	mw := func(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
		return inner(w, r.WithContext(web.SetContextTemplateData(r.Context(), map[string]interface{}{"Public": public})), public)
	}

	return mw
}

func HandleTibiaHtml(w http.ResponseWriter, r *http.Request, isPublicAccess bool) (web.TemplateData, error) {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	tibia := getTibiaData(ag.ID)

	templateData["Tibia"] = tibia

	return templateData, nil
}

func getTibiaData(g int64) map[string]interface{} {
	out := make(map[string]interface{})

	world, _ := GetServerWorld(g, true)
	out["World"] = world

	return out
}
