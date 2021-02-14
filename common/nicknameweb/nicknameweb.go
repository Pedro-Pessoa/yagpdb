package nicknameweb

import (
	"context"
	"database/sql"
	"time"

	"github.com/Pedro-Pessoa/tidbot/analytics"
	"github.com/Pedro-Pessoa/tidbot/bot/eventsystem"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/models"
	"github.com/mediocregopher/radix/v3"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var logger = common.GetPluginLogger(&Plugin{})

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Botnickname",
		SysName:  "botnickname",
		Category: common.PluginCategoryCore,
	}
}

func RegisterPlugin() {
	p := &Plugin{}
	common.RegisterPlugin(p)
}

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLast(p, handleBotMemberUpdate, eventsystem.EventGuildMemberUpdate) // Handle MemberUpdate for nickname changes
}

func handleBotMemberUpdate(evt *eventsystem.EventData) (retry bool, err error) {
	update := evt.GuildMemberUpdate()
	if update.User.ID != common.BotUser.ID {
		return false, nil
	}

	time.Sleep(time.Millisecond * 500)

	var redis bool
	err = common.RedisPool.Do(radix.FlatCmd(&redis, "GET", "cp_nickname_change"))
	if err == nil && redis {
		return false, nil
	}

	conf, err := models.FindCoreConfigG(context.Background(), int64(update.GuildID))
	if err != nil && err != sql.ErrNoRows {
		logger.WithError(err).WithField("guild", update.GuildID).Error("failed fetching core server config")
		return true, err
	}

	if conf == nil {
		conf = &models.CoreConfig{
			GuildID: int64(update.GuildID),
		}
	}

	if conf.BotNickname == update.Member.Nick {
		return false, nil
	}

	conf.BotNickname = update.Member.Nick

	err = conf.Upsert(context.Background(), common.PQ, true, []string{"guild_id"}, boil.Whitelist("bot_nickname"), boil.Infer())
	if err != nil {
		return true, err
	}

	go analytics.RecordActiveUnit(update.GuildID, &Plugin{}, "botnickname_changed")

	return false, nil

}
