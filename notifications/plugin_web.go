package notifications

import (
	"fmt"
	"html/template"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/configstore"
	"github.com/Pedro-Pessoa/tidbot/common/cplogs"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/web"
)

var panelLogKey = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "notifications_settings", FormatString: "Updated server notification settings"})

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../notifications/assets/notifications_general.html", "templates/plugins/notifications_general.html")
	web.AddSidebarItem(web.SidebarCategoryFeeds, &web.SidebarItem{
		Name:   "General",
		NamePT: "Geral",
		URL:    "notifications/general",
		Icon:   "fas fa-bell",
	})

	getHandler := web.RenderHandler(HandleNotificationsGet, "cp_notifications_general")
	postHandler := web.ControllerPostHandler(HandleNotificationsPost, getHandler, Config{})

	web.CPMux.Handle(pat.Get("/notifications/general"), getHandler)
	web.CPMux.Handle(pat.Get("/notifications/general/"), getHandler)

	web.CPMux.Handle(pat.Post("/notifications/general"), postHandler)
	web.CPMux.Handle(pat.Post("/notifications/general/"), postHandler)

	subMux := goji.SubMux()
	subMux.Use(web.NotFound())
}

func HandleNotificationsGet(w http.ResponseWriter, r *http.Request) interface{} {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	formConfig, ok := ctx.Value(common.ContextKeyParsedForm).(*Config)
	if ok {
		templateData["NotifyConfig"] = formConfig
	} else {
		conf, err := GetConfig(activeGuild.ID)
		if err != nil {
			web.CtxLogger(r.Context()).WithError(err).Error("failed retrieving config")
		}

		templateData["NotifyConfig"] = conf
	}

	return templateData
}

func HandleNotificationsPost(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)
	templateData["VisibleURL"] = "/manage/" + discordgo.StrID(activeGuild.ID) + "/notifications/general/"

	newConfig := ctx.Value(common.ContextKeyParsedForm).(*Config)

	newConfig.GuildID = activeGuild.ID

	err := configstore.SQL.SetGuildConfig(ctx, newConfig)
	if err != nil {
		return templateData, nil
	}

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKey))

	return templateData, nil
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "General notifications"
	templateData["WidgetTitlePT"] = "Notificações gerais"
	templateData["SettingsPath"] = "/notifications/general"

	config, err := GetConfig(ag.ID)
	if err != nil {
		return templateData, err
	}

	if templateData["IsPT"] == true {
		const formatPT = `<ul>
		<li>Mensagem de entrada: %s</li>
		<li>Mensagem de entrada (DM): %s</li>
		<li>Mensagem de saída: %s</li>
		<li>Mensagem de troca de tópico: %s</li>
	</ul>`

		templateData["WidgetBodyPT"] = template.HTML(fmt.Sprintf(formatPT,
			web.EnabledDisabledSpanStatusPT(config.JoinServerEnabled), web.EnabledDisabledSpanStatusPT(config.JoinDMEnabled),
			web.EnabledDisabledSpanStatusPT(config.LeaveEnabled), web.EnabledDisabledSpanStatusPT(config.TopicEnabled)))
	} else {
		const format = `<ul>
		<li>Join Server message: %s</li>
		<li>Join DM message: %s</li>
		<li>Leave message: %s</li>
		<li>Topic change message: %s</li>
	</ul>`

		templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format,
			web.EnabledDisabledSpanStatus(config.JoinServerEnabled), web.EnabledDisabledSpanStatus(config.JoinDMEnabled),
			web.EnabledDisabledSpanStatus(config.LeaveEnabled), web.EnabledDisabledSpanStatus(config.TopicEnabled)))
	}

	if config.JoinServerEnabled || config.JoinDMEnabled || config.LeaveEnabled || config.TopicEnabled {
		templateData["WidgetEnabled"] = true
	} else {
		templateData["WidgetDisabled"] = true
	}

	return templateData, nil
}
