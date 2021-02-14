package autorole

import (
	"bytes"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/config"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
	"github.com/vmihailenco/msgpack"
)

var _ = config.RegisterOption("yagpdb.autorole.non_premium_retroactive_assignment", "Wether to enable retroactive assignemnt on non premium guilds", true)

var logger = common.GetPluginLogger(&Plugin{})

func KeyGeneral(guildID int64) string { return "autorole:" + discordgo.StrID(guildID) + ":general" }
func KeyProcessing(guildID int64) string {
	return "autorole:" + discordgo.StrID(guildID) + ":processing"
}

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Autorole",
		SysName:  "autorole",
		Category: common.PluginCategoryMisc,
	}
}

func RegisterPlugin() {
	p := &Plugin{}
	common.RegisterPlugin(p)

	common.GORM.AutoMigrate(&MemberTable{}) // create database table for sticky roles
}

// Autorole general config
type GeneralConfig struct {
	// Autorole
	Role             int64 `json:",string" valid:"role,true"`
	RequiredDuration int
	RequiredRoles    []int64 `valid:"role,true"`
	IgnoreRoles      []int64 `valid:"role,true"`
	OnlyOnJoin       bool

	// Stickyroles
	StickyrolesEnabled bool
	BlacklistedRoles   []int64 `valid:"role,true"`
	WhitelistedRoles   []int64 `valid:"role,true"`
}

// Member table for sticky roles
type MemberTable struct {
	MemberID int64 `gorm:"primary_key"`
	Roles    []byte
}

func serializeValue(v []int64) ([]byte, error) {
	var b bytes.Buffer
	enc := msgpack.NewEncoder(&b)
	err := enc.Encode(v)
	return b.Bytes(), err
}

func deserializeValue(v []byte) ([]int64, error) {
	var out []int64
	if len(v) == 0 {
		return out, nil
	}
	err := msgpack.Unmarshal(v, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func GetGeneralConfig(guildID int64) (*GeneralConfig, error) {
	conf := &GeneralConfig{}
	err := common.GetRedisJson(KeyGeneral(guildID), conf)
	if err != nil {
		logger.WithError(err).WithField("guild", guildID).Error("failed retreiving autorole config")
	}

	return conf, err
}
