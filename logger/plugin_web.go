package logger

import (
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/web"
)

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../logger/assets/logger.html", "templates/plugins/logger.html")

	web.AddSidebarItem(web.SidebarCategoryFeeds, &web.SidebarItem{
		Name:   "Discord Logger",
		NamePT: "Discord Logger",
		URL:    "logger/",
		Icon:   "fas fa-server",
	})

	loggerMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/logger"), loggerMux)
	web.CPMux.Handle(pat.New("/logger/*"), loggerMux)

	loggerMux.Use(web.NotFound())

	loggerGetHandler := web.RenderHandler(HandleLoggerHtml, "cp_logger")

	loggerMux.Handle(pat.Get(""), loggerGetHandler)
	loggerMux.Handle(pat.Get("/"), loggerGetHandler)
}

func HandleLoggerHtml(w http.ResponseWriter, r *http.Request) interface{} {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	logger, ok := ctx.Value(common.ContextKeyParsedForm).(*Logger)
	if ok {
		templateData["Logger"] = logger
	} else {
		conf, err := GetLogger(activeGuild.ID)
		if err != nil {
			web.CtxLogger(r.Context()).WithError(err).Error("failed retrieving config")
		}

		templateData["Logger"] = conf
	}

	return templateData
}
