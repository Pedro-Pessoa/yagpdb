package tibia

import (
	"github.com/jonas747/yagpdb/common"
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
