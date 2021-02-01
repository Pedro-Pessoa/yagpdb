package dutil

import (
	"os"
	"strconv"
	"testing"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var (
	dgo *discordgo.Session // Stores global discordgo session

	envToken   = os.Getenv("DG_TOKEN")            // Token to use when authenticating
	envChannel = parseID(os.Getenv("DG_CHANNEL")) // Channel ID to use for tests
)

func parseID(str string) int64 {
	id, _ := strconv.ParseInt(str, 10, 64)
	return id
}

func init() {
	if envToken == "" {
		return
	}

	if d, err := discordgo.New(envToken); err == nil {
		dgo = d
	}
}

func RequireSession(t *testing.T) bool {
	if dgo == nil || dgo.Token == "" {
		t.Skip("Not logged into discord, skipping...")
		return false
	}

	return true
}

func RequireTestingChannel(t *testing.T) bool {
	if !RequireSession(t) {
		return false
	}

	if envChannel == 0 {
		t.Skip("No testing channel specified, skipping...")
		return false
	}

	return true
}
