package automod_legacy

import (
	"fmt"
	"html/template"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/common/cplogs"
	"github.com/Pedro-Pessoa/tidbot/common/featureflags"
	"github.com/Pedro-Pessoa/tidbot/common/pubsub"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/web"
)

type CtxKey int

const (
	CurrentConfig CtxKey = iota
)

type GeneralForm struct {
	Enabled bool
}

var panelLogKeyUpdatedSettings = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "automod_legacy_settings_updated", FormatString: "Updated legacy automod settings"})

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../automod_legacy/assets/automod_legacy.html", "templates/plugins/automod_legacy.html")

	web.AddSidebarItem(web.SidebarCategoryTools, &web.SidebarItem{
		Name:   "Basic Automoderator",
		NamePT: "Moderador Automático",
		URL:    "automod_legacy",
		Icon:   "fas fa-robot",
	})

	autmodMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/automod_legacy/*"), autmodMux)
	web.CPMux.Handle(pat.New("/automod_legacy"), autmodMux)

	autmodMux.Use(web.NotFound())

	// Alll handlers here require guild channels present
	autmodMux.Use(web.RequireBotMemberMW)
	autmodMux.Use(web.RequirePermMW(discordgo.PermissionManageRoles, discordgo.PermissionKickMembers, discordgo.PermissionBanMembers, discordgo.PermissionManageMessages))

	getHandler := web.RenderHandler(HandleAutomod, "cp_automod_legacy")

	autmodMux.Handle(pat.Get("/"), getHandler)
	autmodMux.Handle(pat.Get(""), getHandler)

	// Post handlers
	autmodMux.Handle(pat.Post("/"), ExtraPostMW(web.SimpleConfigSaverHandler(Config{}, getHandler, panelLogKeyUpdatedSettings)))
	autmodMux.Handle(pat.Post(""), ExtraPostMW(web.SimpleConfigSaverHandler(Config{}, getHandler, panelLogKeyUpdatedSettings)))
}

func HandleAutomod(w http.ResponseWriter, r *http.Request) interface{} {
	g, templateData := web.GetBaseCPContextData(r.Context())

	config, err := GetConfig(g.ID)
	web.CheckErr(templateData, err, "Failed retrieving rules", web.CtxLogger(r.Context()).Error)

	templateData["AutomodConfig"] = config
	templateData["VisibleURL"] = "/manage/" + discordgo.StrID(g.ID) + "/automod_legacy/"

	return templateData
}

// Invalidates the cache
func ExtraPostMW(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		activeGuild, _ := web.GetBaseCPContextData(r.Context())
		_ = pubsub.Publish("update_automod_legacy_rules", activeGuild.ID, nil)
		featureflags.MarkGuildDirty(activeGuild.ID)
		inner.ServeHTTP(w, r)
	}

	return http.HandlerFunc(mw)
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	g, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "Basic Automod"
	templateData["WidgetTitlePT"] = "Moderador Automático"
	templateData["SettingsPath"] = "/automod_legacy"

	config, err := GetConfig(g.ID)
	if err != nil {
		return templateData, err
	}

	if config.Enabled {
		templateData["WidgetEnabled"] = true
	} else {
		templateData["WidgetDisabled"] = true
	}

	if templateData["IsPT"] == true {
		const formatPT = `<ul>
		<li>Modo lento: %s</li>
		<li>Spam de menção: %s</li>
		<li>Convites de servidores: %s</li>
		<li>Quaisquer links: %s</li>
		<li>Palavras banidas: %s</li>
		<li>Sites banidos: %s</li>
	</ul>`

		slowmodePT := web.EnabledDisabledSpanStatusPT(config.Spam.Enabled)
		massMentionPT := web.EnabledDisabledSpanStatusPT(config.Mention.Enabled)
		invitesPT := web.EnabledDisabledSpanStatusPT(config.Invite.Enabled)
		linksPT := web.EnabledDisabledSpanStatusPT(config.Links.Enabled)
		wordsPT := web.EnabledDisabledSpanStatusPT(config.Words.Enabled)
		sitesPT := web.EnabledDisabledSpanStatusPT(config.Sites.Enabled)

		templateData["WidgetBodyPT"] = template.HTML(fmt.Sprintf(formatPT, slowmodePT,
			massMentionPT, invitesPT, linksPT, wordsPT, sitesPT))
	} else {
		const format = `<ul>
	<li>Slowmode: %s</li>
	<li>Mass mention: %s</li>
	<li>Server invites: %s</li>
	<li>Any links: %s</li>
	<li>Banned words: %s</li>
	<li>Banned websites: %s</li>
	</ul>`

		slowmode := web.EnabledDisabledSpanStatus(config.Spam.Enabled)
		massMention := web.EnabledDisabledSpanStatus(config.Mention.Enabled)
		invites := web.EnabledDisabledSpanStatus(config.Invite.Enabled)
		links := web.EnabledDisabledSpanStatus(config.Links.Enabled)
		words := web.EnabledDisabledSpanStatus(config.Words.Enabled)
		sites := web.EnabledDisabledSpanStatus(config.Sites.Enabled)

		templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, slowmode,
			massMention, invites, links, words, sites))
	}

	return templateData, nil
}
