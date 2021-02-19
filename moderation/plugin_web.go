package moderation

import (
	"fmt"
	"html/template"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/cplogs"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/Pedro-Pessoa/tidbot/web"
)

var (
	panelLogKeyUpdatedSettings = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "moderation_settings_updated", FormatString: "Updated moderation config"})
	panelLogKeyClearWarnings   = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "moderation_warnings_cleared", FormatString: "Cleared %d moderation user warnings"})
)

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../moderation/assets/moderation.html", "templates/plugins/moderation.html")

	web.AddSidebarItem(web.SidebarCategoryTools, &web.SidebarItem{
		Name:   "Moderation",
		NamePT: "Moderação",
		URL:    "moderation",
		Icon:   "fas fa-gavel",
	})

	subMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/moderation"), subMux)
	web.CPMux.Handle(pat.New("/moderation/*"), subMux)

	subMux.Use(web.RequireBotMemberMW) // need the bot's role
	subMux.Use(web.RequirePermMW(discordgo.PermissionManageRoles, discordgo.PermissionKickMembers, discordgo.PermissionBanMembers, discordgo.PermissionManageMessages, discordgo.PermissionEmbedLinks))
	subMux.Use(web.NotFound())

	getHandler := web.ControllerHandler(HandleModeration, "cp_moderation")
	postHandler := web.ControllerPostHandler(HandlePostModeration, getHandler, Config{})
	clearServerWarnings := web.ControllerPostHandler(HandleClearServerWarnings, getHandler, nil)

	subMux.Handle(pat.Get(""), getHandler)
	subMux.Handle(pat.Get("/"), getHandler)
	subMux.Handle(pat.Post(""), postHandler)
	subMux.Handle(pat.Post("/"), postHandler)
	subMux.Handle(pat.Post("/clear_server_warnings"), clearServerWarnings)
}

// HandleModeration servers the moderation page itself
func HandleModeration(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	activeGuild, templateData := web.GetBaseCPContextData(r.Context())

	templateData["DefaultDMMessage"] = DefaultDMMessage

	if _, ok := templateData["ModConfig"]; !ok {
		config, err := GetConfig(activeGuild.ID)
		if err != nil {
			return templateData, err
		}
		templateData["ModConfig"] = config
	}

	return templateData, nil
}

// HandlePostModeration update the settings
func HandlePostModeration(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)
	templateData["VisibleURL"] = "/manage/" + discordgo.StrID(activeGuild.ID) + "/moderation/"

	newConfig := ctx.Value(common.ContextKeyParsedForm).(*Config)
	newConfig.DefaultMuteDuration.Valid = true
	newConfig.DefaultBanDeleteDays.Valid = true
	newConfig.DefaultLockdownDuration.Valid = true
	templateData["ModConfig"] = newConfig

	err := newConfig.Save(activeGuild.ID)

	templateData["DefaultDMMessage"] = DefaultDMMessage

	if err == nil {
		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyUpdatedSettings))
	}

	return templateData, err
}

// Clear all server warnigns
func HandleClearServerWarnings(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)
	templateData["VisibleURL"] = "/manage/" + discordgo.StrID(activeGuild.ID) + "/moderation/"

	rows := common.GORM.Where("guild_id = ?", activeGuild.ID).Delete(WarningModel{}).RowsAffected
	templateData.AddAlerts(web.SucessAlert("Deleted ", rows, " warnings!"))
	templateData["DefaultDMMessage"] = DefaultDMMessage

	if rows > 0 {
		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyClearWarnings, &cplogs.Param{Type: cplogs.ParamTypeInt, Value: rows}))
	}

	return templateData, nil
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	activeGuild, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "Moderation"
	templateData["WidgetTitlePT"] = "Moderação"
	templateData["SettingsPath"] = "/moderation"

	config, err := GetConfig(activeGuild.ID)
	if err != nil {
		return templateData, err
	}

	if templateData["IsPT"] == true {
		const formatPT = `<ul>
		<li>Comando reportar: %s</li>
		<li>Comando limpar: %s</li>
		<li>Comando dar/tirar cargo: %s</li>
		<li>Comando expulsar: %s</li>
		<li>Comando banir: %s</li>
		<li>Comando silenciar: %s</li>
		<li>Comando avisar: %s</li>
		<li>Comando lockdown: %s</li>
		<li>Comando modo lento: %s</li>
	</ul>`

		templateData["WidgetBodyPT"] = template.HTML(fmt.Sprintf(formatPT, web.EnabledDisabledSpanStatusPT(config.ReportEnabled),
			web.EnabledDisabledSpanStatusPT(config.CleanEnabled), web.EnabledDisabledSpanStatusPT(config.GiveRoleCmdEnabled),
			web.EnabledDisabledSpanStatusPT(config.KickEnabled), web.EnabledDisabledSpanStatusPT(config.BanEnabled),
			web.EnabledDisabledSpanStatusPT(config.MuteEnabled), web.EnabledDisabledSpanStatusPT(config.WarnCommandsEnabled),
			web.EnabledDisabledSpanStatusPT(config.LockdownCmdEnabled), web.EnabledDisabledSpanStatusPT(config.SlowmodeCommandEnabled)))
	} else {
		const format = `<ul>
		<li>Report command: %s</li>
		<li>Clean command: %s</li>
		<li>Giverole/Takerole commands: %s</li>
		<li>Kick command: %s</li>
		<li>Ban command: %s</li>
		<li>Mute/Unmute commands: %s</li>
		<li>Warning commands: %s</li>
		<li>Lockdown commands: %s</li>
		<li>Slowmode command: %s</li>
	</ul>`

		templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, web.EnabledDisabledSpanStatus(config.ReportEnabled),
			web.EnabledDisabledSpanStatus(config.CleanEnabled), web.EnabledDisabledSpanStatus(config.GiveRoleCmdEnabled),
			web.EnabledDisabledSpanStatus(config.KickEnabled), web.EnabledDisabledSpanStatus(config.BanEnabled),
			web.EnabledDisabledSpanStatus(config.MuteEnabled), web.EnabledDisabledSpanStatus(config.WarnCommandsEnabled),
			web.EnabledDisabledSpanStatus(config.LockdownCmdEnabled), web.EnabledDisabledSpanStatus(config.SlowmodeCommandEnabled)))
	}

	if config.ReportEnabled || config.CleanEnabled || config.GiveRoleCmdEnabled || config.ActionChannel != "" ||
		config.MuteEnabled || config.KickEnabled || config.BanEnabled || config.WarnCommandsEnabled || config.LockdownCmdEnabled || config.SlowmodeCommandEnabled {
		templateData["WidgetEnabled"] = true
	} else {
		templateData["WidgetDisabled"] = true
	}

	return templateData, nil
}
