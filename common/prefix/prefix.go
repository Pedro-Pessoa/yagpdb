package prefix

import (
	"strings"

	"github.com/mediocregopher/radix/v3"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/featureflags"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

func GetCommandPrefixRedis(guild int64) (string, error) {
	var prefix string
	err := common.RedisPool.Do(radix.Cmd(&prefix, "GET", "command_prefix:"+discordgo.StrID(guild)))
	if err == nil && prefix == "" {
		prefix = DefaultCommandPrefix()
	}

	return strings.TrimSpace(prefix), err
}

func DefaultCommandPrefix() string {
	defaultPrefix := "!"
	if common.Testing {
		defaultPrefix = "("
	}

	return defaultPrefix
}

func GetPrefixIgnoreError(guild int64) string {
	prefix := DefaultCommandPrefix()
	if featureflags.GuildHasFlagOrLogError(guild, "commands_has_custom_prefix") {
		prefix, _ = GetCommandPrefixRedis(guild)
	}

	return strings.TrimSpace(prefix)
}
