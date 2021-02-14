package tidstatus

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var Command = &commands.TIDCommand{
	Cooldown:    5,
	CmdCategory: commands.CategoryDebug,
	Name:        "Tidstatus",
	Aliases:     []string{"status"},
	Description: "Shows TidBot status, version, uptime, memory stats, and so on",
	RunInDM:     true,
	RunFunc:     cmdFuncYagStatus,
}

var logger = common.GetFixedPrefixLogger("yagstatuc_cmd")

func cmdFuncYagStatus(data *dcmd.Data) (interface{}, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	servers, _ := common.GetJoinedServerCount()
	sysMem, err := mem.VirtualMemory()
	sysMemStats := ""
	if err == nil {
		sysMemStats = fmt.Sprintf("%dMB (%.0f%%), %dMB", sysMem.Used/1000000, sysMem.UsedPercent, sysMem.Total/1000000)
	} else {
		sysMemStats = "Failed collecting mem stats"
		logger.WithError(err).Error("Failed collecting memory stats")
	}

	sysLoad, err := load.Avg()
	sysLoadStats := ""
	if err == nil {
		sysLoadStats = fmt.Sprintf("%.2f, %.2f, %.2f", sysLoad.Load1, sysLoad.Load5, sysLoad.Load15)
	} else {
		sysLoadStats = "Failed collecting"
		logger.WithError(err).Error("Failed collecting load stats")
	}

	uptime := time.Since(bot.Started)
	allocated := float64(memStats.Alloc) / 1000000
	numGoroutines := runtime.NumGoroutine()
	botUser := common.BotUser

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    botUser.Username,
			IconURL: discordgo.EndpointUserAvatar(botUser.ID, botUser.Avatar),
		},
		Title: "Tid Bot Status, build version " + common.VERSION,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Servers", Value: fmt.Sprint(servers), Inline: true},
			{Name: "Go Version", Value: runtime.Version(), Inline: true},
			{Name: "API Version", Value: discordgo.APIVersion, Inline: true},
			{Name: "Uptime", Value: common.HumanizeDuration(common.DurationPrecisionSeconds, uptime), Inline: true},
			{Name: "Goroutines", Value: fmt.Sprint(numGoroutines), Inline: true},
			{Name: "GC Pause Fraction", Value: fmt.Sprintf("%.3f%%", memStats.GCCPUFraction*100), Inline: true},
			{Name: "Process Mem (alloc, sys, freed)", Value: fmt.Sprintf("%.1fMB, %.1fMB, %.1fMB", float64(memStats.Alloc)/1000000, float64(memStats.Sys)/1000000, (float64(memStats.TotalAlloc)/1000000)-allocated), Inline: true},
			{Name: "System Mem (used, total)", Value: sysMemStats, Inline: true},
			{Name: "System Load (1, 5, 15)", Value: sysLoadStats, Inline: true},
			{Name: "Master version", Value: common.CurrentVersion, Inline: true},
		},
	}

	for _, v := range common.Plugins {
		if cast, ok := v.(PluginStatus); ok {
			started := time.Now()
			name, val := cast.Status()
			if name == "" || val == "" {
				continue
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: v.PluginInfo().Name + ": " + name, Value: val, Inline: true})
			elapsed := time.Since(started)
			logger.Println("Took ", elapsed.Seconds(), " to gather stats from ", v.PluginInfo().Name)
		}
	}

	return embed, nil
	// return &commandsystem.FallbackEmebd{embed}, nil
}

type PluginStatus interface {
	Status() (string, string)
}
