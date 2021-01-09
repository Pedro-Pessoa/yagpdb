package aylien

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
)

var (
	ErrNoMessages = errors.New("Failed finding any messages to analyze")

	logger = common.GetPluginLogger(&Plugin{})
)

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "AYLIEN",
		SysName:  "aylien",
		Category: common.PluginCategoryMisc,
	}
}

func RegisterPlugin() {
	plugin := &Plugin{}

	common.RegisterPlugin(plugin)
}

var _ commands.CommandProvider = (*Plugin)(nil)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, &commands.YAGCommand{
		CmdCategory: commands.CategoryFun,
		Cooldown:    5,
		Name:        "Sentiment",
		Aliases:     []string{"sent"},
		Description: "Does sentiment analysis on a message or your last 5 messages longer than 3 words",
		Arguments: []*dcmd.ArgDef{
			{Name: "text", Type: dcmd.String},
		},
		RunFunc: func(cmd *dcmd.Data) (interface{}, error) {
			toAnalyze := make([]string, 0)

			if cmd.Args[0].Value != nil { // a message was provided
				toAnalyze = append(toAnalyze, cmd.Args[0].Str())
			} else { // Get the message to analyze
				msgs, err := bot.GetMessages(cmd.CS.ID, 100, false)
				if err != nil {
					return "", err
				}

				if len(msgs) < 1 {
					return ErrNoMessages, ErrNoMessages
				}

				for _, m := range msgs { // filter out our own and longer than 3 words
					if m.Author.ID == cmd.Msg.Author.ID {
						if len(strings.Fields(m.ContentWithMentionsReplaced())) > 3 {
							toAnalyze = append(toAnalyze, m.Content)
							if len(toAnalyze) >= 5 {
								break
							}
						}
					}
				}

				if len(toAnalyze) < 1 {
					return ErrNoMessages, ErrNoMessages
				}
			}

			rand.Seed(time.Now().UnixNano())

			out := fmt.Sprintf("**Sentiment analysis on %d message(s):**\n", len(toAnalyze))
			for _, resp := range toAnalyze {
				out += fmt.Sprintf("*%s*\nPolarity: **%s** *(Confidence: %s%%)* Subjectivity: **%s** *(Confidence: %s%%)*\n\n", resp, RandPolarity(), RandConfidencePol(), RandSubjective(), RandConfidenceSub())
			}

			return out, nil
		},
	},
		&commands.YAGCommand{
			Cooldown:    2,
			CmdCategory: commands.CategoryFun,
			Name:        "8Ball",
			Description: "Wisdom",
			Arguments: []*dcmd.ArgDef{
				{Name: "What to ask", Type: dcmd.String},
			},
			RequiredArgs: 1,
			RunFunc: func(cmd *dcmd.Data) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				return fmt.Sprintf("**8Ball**:\n_%s_\n%s", cmd.Args[0].Str(), Rand8Ball()), nil
			},
		},
	)
}

func RandPolarity() string {
	i := rand.Intn(3)
	switch i {
	case 0:
		return "neutral"
	case 1:
		return "positive"
	default:
		return "negative"
	}
}

func RandConfidencePol() string {
	i := rand.Intn(10001)
	idiv := float64(i) / 100

	return fmt.Sprintf("%v", idiv)
}

func RandSubjective() string {
	i := rand.Intn(2)
	switch i {
	case 0:
		return "subjective"
	default:
		return "objective"
	}
}

func RandConfidenceSub() string {
	i := rand.Intn(10001)
	idiv := float64(i) / 100

	return fmt.Sprintf("%v", idiv)
}

func Rand8Ball() string {
	i := rand.Intn(6)
	switch i {
	case 0:
		return "Yes"
	case 1:
		return "No"
	case 2:
		return "Definitively not"
	case 3:
		return "Not likely"
	case 4:
		return "Most likely"
	default:
		return "Without a doubt"
	}
}
