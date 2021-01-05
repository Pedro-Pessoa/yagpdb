package tibia

import (
	"github.com/jinzhu/gorm"
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

	common.GORM.AutoMigrate(&TibiaFlags{}, &TibiaTracking{}, &ScanTable{})

	table := ScanTable{}
	err := common.GORM.Where(&table).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if alreadySet {
		if table.RunScan {
			store := New()
			store.TrackingController()
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
