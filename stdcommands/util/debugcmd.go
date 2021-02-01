package util

import (
	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
)

func IsExecedByCC(data *dcmd.Data) bool {
	if v := data.Context().Value(commands.CtxKeyExecutedByCC); v != nil {
		if cast, _ := v.(bool); cast {
			return true
		}
	}

	return false
}

func RequireBotAdmin(inner dcmd.RunFunc) dcmd.RunFunc {
	return func(data *dcmd.Data) (interface{}, error) {
		if IsExecedByCC(data) {
			return "", nil
		}

		if admin, err := bot.IsBotAdmin(data.Msg.Author.ID); admin && err == nil {
			return inner(data)
		}

		return "", nil
	}
}

func RequireOwner(inner dcmd.RunFunc) dcmd.RunFunc {
	return func(data *dcmd.Data) (interface{}, error) {
		if IsExecedByCC(data) {
			return "", nil
		}

		if common.IsOwner(data.Msg.Author.ID) {
			return inner(data)
		}

		return "", nil
	}
}
