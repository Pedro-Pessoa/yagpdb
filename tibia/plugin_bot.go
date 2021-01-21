package tibia

import (
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/mqueue"
)

var logger = common.GetPluginLogger(&Plugin{})

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Tibia",
		SysName:  "tibia",
		Category: common.PluginCategoryTibia,
	}
}

func RegisterPlugin() {
	plugin := &Plugin{}

	common.RegisterPlugin(plugin)

	common.GORM.AutoMigrate(&TibiaFlags{}, &TibiaTracking{}, &ScanTable{}, &InnerNewsStruct{}, &NewsTable{})

	mqueue.RegisterSource("tibia", plugin)

	table := ScanTable{}
	err := common.GORM.Where(&table).First(&table).Error
	if err == nil {
		if table.RunScan {
			store := New()
			store.trackingController()
		}
	}

	newsTable := NewsTable{}
	err = common.GORM.Where(&newsTable).First(&newsTable).Error
	if err == nil {
		if newsTable.RunScan {
			newsController()
		}
	}
}

func (p *Plugin) DisableFeed(elem *mqueue.QueuedElement, PlaceHolder error) {
	table := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: elem.Guild}).First(&table).Error
	if err != nil {
		logger.Errorf("Error getting track on disable feed: %#v", err)
	} else {
		table.SendDeaths = false
		table.SendUpdates = false
		err = common.GORM.Save(&table).Error
		if err != nil {
			logger.Errorf("Error saving track on disable feed: %#v", err)
		}
	}

	newstable := InnerNewsStruct{}
	err = common.GORM.Where(&InnerNewsStruct{GuildID: elem.Guild}).First(&newstable).Error
	if err != nil {
		logger.Errorf("Error getting news on disable feed: %#v", err)
	} else {
		newstable.RunNews = false
		err = common.GORM.Save(&newstable).Error
		if err != nil {
			logger.Errorf("Error saving news on disable feed: %#v", err)
		}
	}
}

var _ mqueue.PluginWithSourceDisabler = (*Plugin)(nil)

// func (p *Plugin) AddCommands() {
//	commands.AddRootCommands(p, RespCommands...)
// }

// func (p *Plugin) BotInit() {
// 	//	scheduledevents2.RegisterHandler("resp_handler", RespHandler{}, HandleRespawns)

// 	eventsystem.AddHandlerAsyncLast(p, ConvertSystem(), eventsystem.EventMessageDelete)
// }

// func ConvertSystem() eventsystem.HandlerFunc {
// 	return func(evt *eventsystem.EventData) (retry bool, err error) {
// 		HandleMessageDelete(evt)
// 		return false, nil
// 	}
// }
