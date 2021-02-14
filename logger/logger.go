package logger

import (
	"context"
	"strconv"

	"emperror.dev/errors"
	"github.com/lib/pq"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/configstore"
	"github.com/Pedro-Pessoa/tidbot/common/featureflags"
)

var logger = common.GetPluginLogger(&Plugin{})

// Plugin structure for discord looger
type Plugin struct{}

// Plugin info for discord logger
func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Logger",
		SysName:  "logger",
		Category: common.PluginCategoryMisc,
	}
}

// Registers discord logger plugin
func RegisterPlugin() {
	p := &Plugin{}
	common.RegisterPlugin(p)

	common.GORM.AutoMigrate(&Logger{})
	configstore.RegisterConfig(configstore.SQL, &Logger{})
}

type Logger struct {
	configstore.GuildConfigModel

	// Global settings
	LoggerEnabled            bool          `json:"logger_enabled" schema:"logger_enabled"`
	DefaultLoggerChannel     string        `json:"default_logger_channel" schema:"default_logger_channel" valid:"channel,true"`
	IgnoredLoggerChannels    pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`
	RequiredLoggerChannels   pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`
	RequiredLoggerCategories pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`
	IgnoredLoggerCategories  pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`
	RequiredLoggerRoles      pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`
	IgnoredLoggerRoles       pq.Int64Array `gorm:"type:bigint[]" valid:"channel,true"`

	// Server events
	ServerEventsEnabled    bool   `json:"server_events_enabled" schema:"server_events_enabled"`
	ServerLoggerChannel    string `json:"server_logger_channel" schema:"server_logger_channel" valid:"channel,true"`
	ChannelCreationEnabled bool   `json:"channel_creation_enabled" schema:"channel_creation_enabled"`
	ChannelUpdatedEnabled  bool   `json:"channel_updated_enabled" schema:"channel_updated_enabled"`
	ChannelDeletionEnbaled bool   `json:"channe_deletion_enabled" schema:"channe_deletion_enabled"`
	RoleCreationEnbaled    bool   `json:"role_creation_enabled" schema:"role_creation_enabled"`
	RoleUpdateEnbaled      bool   `json:"role_updated_enabled" schema:"role_updated_enabled"`
	RoleDeletionEnbaled    bool   `json:"role_deletion_enabled" schema:"role_deletion_enabled"`
	ServerUpdateEnbaled    bool   `json:"server_update_enabled" schema:"server_update_enabled"`
	EmojiUpdateEnbaled     bool   `json:"emoji_update_enabled" schema:"emoji_update_enabled"`
	ServerLoggerMSG        string `json:"server_logger_msg" schema:"server_logger_msg" valid:"template,2000"`

	// Messages
	MessageEventsEnabled   bool   `json:"message_events_enabled" schema:"message_events_enabled"`
	MessageLoggerChannel   string `json:"message_logger_channel" schema:"message_logger_channel" valid:"channel,true"`
	DeletedMessagesEnabled bool   `json:"deleted_messages_enabled" schema:"deleted_messages_enabled"`
	EditedMessagesEnabled  bool   `json:"edited_messages_enabled" schema:"edited_messages_enabled"`
	PurgedMessagesEnabled  bool   `json:"purged_messages_enabled" schema:"purged_messages_enabled"`
	MessageLoggerMSG       string `json:"message_logger_msg" schema:"message_logger_msg" valid:"template,2000"`

	// Members
	MemberEventsEnabled       bool   `json:"member_events_enabled" schema:"member_events_enabled"`
	MemberLoggerChannel       string `json:"member_logger_channel" schema:"member_logger_channel" valid:"channel,true"`
	MemberRoleUpdateEnabled   bool   `json:"member_role_update_enabled" schema:"member_role_update_enabled"`
	MemberNameUpdateEnabled   bool   `json:"member_name_update_enabled" schema:"member_name_update_enabled"`
	MemberAvatarUpdateEnabled bool   `json:"member_avatar_update_enabled" schema:"member_avatar_update_enabled"`
	MemberLoggerMSG           string `json:"member_logger_msg" schema:"member_logger_msg" valid:"template,2000"`

	// Voices
	VoiceEventsEnabled bool   `json:"voice_events_enabled" schema:"voice_events_enabled"`
	VoiceLoggerChannel string `json:"voice_logger_channel" schema:"voice_logger_channel" valid:"channel,true"`
	JoinVoiceEnabled   bool   `json:"join_voice_enabled" schema:"join_voice_enabled"`
	SwapVoiceEnabled   bool   `json:"swap_voice_update_enabled" schema:"swap_voice_update_enabled"`
	LeaveVoiceEnabled  bool   `json:"leave_voice_update_enabled" schema:"leave_voice_update_enabled"`
	VoiceLoggerMSG     string `json:"voice_logger_msg" schema:"voice_logger_msg" valid:"template,2000"`
}

func (c *Logger) DefaultLoggerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.DefaultLoggerChannel, 10, 64)
	return
}

func (c *Logger) ServerLoggerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.ServerLoggerChannel, 10, 64)
	return
}

func (c *Logger) MessageLoggerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.MessageLoggerChannel, 10, 64)
	return
}

func (c *Logger) MemberLoggerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.MemberLoggerChannel, 10, 64)
	return
}

func (c *Logger) VoiceLoggerChannelInt() (i int64) {
	i, _ = strconv.ParseInt(c.VoiceLoggerChannel, 10, 64)
	return
}

func (c *Logger) GetName() string {
	return "discord_logger"
}

func (c *Logger) TableName() string {
	return "discord_logger_configs"
}

func GetLogger(guildID int64) (*Logger, error) {
	var conf Logger
	err := configstore.Cached.GetGuildConfig(context.Background(), guildID, &conf)
	if err == nil {
		return &conf, nil
	}

	if err == configstore.ErrNotFound {
		return &Logger{
			ServerLoggerMSG:  DefaultMessage,
			MessageLoggerMSG: DefaultMessage,
			MemberLoggerMSG:  DefaultMessage,
			VoiceLoggerMSG:   DefaultMessage,
		}, nil
	}

	return nil, errors.WithStackIf(err)
}

var _ featureflags.PluginWithFeatureFlags = (*Plugin)(nil)

const (
	featureFlagEnabled = "discord_logger_enabled"
)

func (p *Plugin) UpdateFeatureFlags(guildID int64) ([]string, error) {
	logger, err := GetLogger(guildID)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	if !logger.LoggerEnabled && !logger.ServerEventsEnabled && !logger.MessageEventsEnabled && !logger.MemberEventsEnabled && !logger.VoiceEventsEnabled {
		return nil, nil
	}

	return []string{featureFlagEnabled}, nil
}

func (p *Plugin) AllFeatureFlags() []string {
	return []string{
		featureFlagEnabled, // set if there is atleast one ruleset enabled with a rule in it
	}
}

const DefaultMessage = `{{$channelTypes := dict
	0 "Text Channel"
	2 "Voice Channel"
	4 "Guild Category"
	5 "Guild News"
	6 "Guild Store"
}}

{{$channelType := $channelTypes.Get .EventChannel.Type}}

{{$channelMap := sdict
	"Channel Created" (printf "**Channel name**: %s\n**Channel ID**: %d\n**Channel Link**: <#%d>\n**Channel Type**: %s\n**Channel Topic**: %s\n**NSFW**: %v" .EventChannel.Name .EventChannel.ID .EventChannel.ID $channelType .EventChannel.Topic .EventChannel.NSFW)
	"Channel Updated" (printf "%s\n**Channel ID**: %d\n**Channel Link**: <#%d>" .Changes .EventChannel.ID .EventChannel.ID)
	"Channel Deleted" (printf "Channel **%s** was deleted on category <#%d>" .EventChannel.Name, .EventChannel.ParentID)
}}

{{$roleMap := sdict
	"Role Created" (printf "**Role Name**: %s\n**Role ID**: %d\n**Role Mention**: <@&%d>\n**Manged**: %v\n**Mentionbale**: %v\n**Hoisted**: %v\n**Color**: %d" .Role.Name .Role.ID .Role.ID .Role.Managed .Role.Mentionable .Role.Hoisted .Role.Color)
	"Role Updated" (printf "%s\n**Role ID**: %d\n**Role Mention**: <@&%d>" .Changes .Role.ID .Role.ID)
	"Role Deleted" (printf "Role **%s** was deleted." .Role.Name)
}}

{{$eventsMap := sdict
	"ChannelEvent" $channelMap
	"RoleEvent" $roleMap
}}

{{$descr := ""}}
{{with $eventsMap.Get .EventType}}
	{{with .Get $.EventName}}
		{{$descr = .}}
	{{end}}
{{end}}

{{$embed := cembed "title" .EventName "description" $descr "color" .DefaultColor "timestamp" currentTime}}

{{sendMessage nil $embed}}
`
