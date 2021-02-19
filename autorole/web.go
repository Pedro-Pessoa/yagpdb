package autorole

import (
	"fmt"
	"html"
	"html/template"
	"net/http"

	"emperror.dev/errors"
	"github.com/mediocregopher/radix/v3"
	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/cplogs"
	"github.com/Pedro-Pessoa/tidbot/common/pubsub"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/web"
)

type Form struct {
	GeneralConfig `valid:"traverse"`
}

var _ web.SimpleConfigSaver = (*Form)(nil)

var (
	panelLogKeyUpdatedSettings = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "autorole_settings_updated", FormatString: "Updated autorole settings"})
	panelLogKeyStartedFullScan = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "autorole_full_scan", FormatString: "Started full retroactive autorole scan"})
)

func (f Form) Save(guildID int64) error {
	_ = pubsub.Publish("autorole_stop_processing", guildID, nil)

	err := common.SetRedisJson(KeyGeneral(guildID), f.GeneralConfig)
	if err != nil {
		return err
	}

	return nil
}

func (f Form) Name() string {
	return "Autorole"
}

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../autorole/assets/autorole.html", "templates/plugins/autorole.html")

	web.AddSidebarItem(web.SidebarCategoryTools, &web.SidebarItem{
		Name:   "Autorole",
		NamePT: "Cargo Automático",
		URL:    "autorole",
		Icon:   "fas fa-user-plus",
	})

	muxer := goji.SubMux()

	web.CPMux.Handle(pat.New("/autorole"), muxer)
	web.CPMux.Handle(pat.New("/autorole/*"), muxer)

	muxer.Use(web.RequireBotMemberMW) // need the bot's role
	muxer.Use(web.RequirePermMW(discordgo.PermissionManageRoles))
	muxer.Use(web.NotFound())

	getHandler := web.RenderHandler(handleGetAutoroleMainPage, "cp_autorole")

	muxer.Handle(pat.Get(""), getHandler)
	muxer.Handle(pat.Get("/"), getHandler)

	muxer.Handle(pat.Post("/fullscan"), web.ControllerPostHandler(handlePostFullScan, getHandler, nil))

	muxer.Handle(pat.Post(""), web.SimpleConfigSaverHandler(Form{}, getHandler, panelLogKeyUpdatedSettings))
	muxer.Handle(pat.Post("/"), web.SimpleConfigSaverHandler(Form{}, getHandler, panelLogKeyUpdatedSettings))
}

func handleGetAutoroleMainPage(w http.ResponseWriter, r *http.Request) interface{} {
	ctx := r.Context()
	activeGuild, tmpl := web.GetBaseCPContextData(ctx)

	general, err := GetGeneralConfig(activeGuild.ID)
	web.CheckErr(tmpl, err, "Failed retrieving general config (contact support)", web.CtxLogger(r.Context()).Error)
	tmpl["Autorole"] = general

	var proc int
	_ = common.RedisPool.Do(radix.Cmd(&proc, "GET", KeyProcessing(activeGuild.ID)))
	tmpl["Processing"] = proc
	tmpl["ProcessingETA"] = int(proc / 60)

	fullScanActive := WorkingOnFullScan(activeGuild.ID)
	tmpl["FullScanActive"] = fullScanActive

	return tmpl

}

func handlePostFullScan(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, tmpl := web.GetBaseCPContextData(ctx)

	err := botRestPostFullScan(activeGuild.ID)
	if err != nil {
		if err == ErrAlreadyProcessingFullGuild {
			return tmpl.AddAlerts(web.ErrorAlert("Already processing, please wait.")), nil
		}

		return tmpl, errors.WithMessage(err, "botrest")
	}

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyStartedFullScan))

	return tmpl, nil
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "Autorole"
	templateData["WidgetTitlePT"] = "Cargo Automático"
	templateData["SettingsPath"] = "/autorole"

	general, err := GetGeneralConfig(ag.ID)
	if err != nil {
		return templateData, err
	}

	if templateData["IsPT"] == true {
		var enabledDisabledPT string
		autoroleRolePT := "nenhum"

		if role := ag.Role(general.Role); role != nil {
			templateData["WidgetEnabled"] = true
			enabledDisabledPT = web.EnabledDisabledSpanStatusPT(true)
			autoroleRolePT = html.EscapeString(role.Name)
		} else {
			templateData["WidgetDisabled"] = true
			enabledDisabledPT = web.EnabledDisabledSpanStatusPT(false)
		}

		const formatPT = `<ul>
	<li>Status do cargo automático: %s</li>
	<li>Cargo: <code>%s</code></li>
	</ul>`

		templateData["WidgetBodyPT"] = template.HTML(fmt.Sprintf(formatPT, enabledDisabledPT, autoroleRolePT))
	} else {
		var enabledDisabled string
		autoroleRole := "none"

		if role := ag.Role(general.Role); role != nil {
			templateData["WidgetEnabled"] = true
			enabledDisabled = web.EnabledDisabledSpanStatus(true)
			autoroleRole = html.EscapeString(role.Name)
		} else {
			templateData["WidgetDisabled"] = true
			enabledDisabled = web.EnabledDisabledSpanStatus(false)
		}

		const format = `<ul>
	<li>Autorole status: %s</li>
	<li>Autorole role: <code>%s</code></li>
	</ul>`

		templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, enabledDisabled, autoroleRole))
	}

	return templateData, nil
}
