package reminders

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jinzhu/gorm"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/commands"
	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/common/scheduledevents2"
	seventsmodels "github.com/Pedro-Pessoa/tidbot/common/scheduledevents2/models"
	"github.com/Pedro-Pessoa/tidbot/pkgs/dcmd"
	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

var (
	logger = common.GetPluginLogger(&Plugin{})

	_ bot.BotInitHandler       = (*Plugin)(nil)
	_ commands.CommandProvider = (*Plugin)(nil)
)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, cmds...)
}

func (p *Plugin) BotInit() {
	// scheduledevents.RegisterEventHandler("reminders_check_user", checkUserEvtHandlerLegacy)
	scheduledevents2.RegisterHandler("reminders_check_user", int64(0), checkUserScheduledEvent)
	scheduledevents2.RegisterLegacyMigrater("reminders_check_user", migrateLegacyScheduledEvents)
}

// Reminder management commands
var cmds = []*commands.TIDCommand{
	{
		CmdCategory:  commands.CategoryTool,
		Name:         "Remindme",
		Description:  "Schedules a reminder, example: 'remindme 1h30min are you alive still?'",
		Aliases:      []string{"remind", "reminder"},
		RequiredArgs: 2,
		Cooldown:     5,
		Arguments: []*dcmd.ArgDef{
			{Name: "Time", Type: &commands.DurationArg{}},
			{Name: "Message", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			currentReminders, _ := GetUserReminders(parsed.Msg.Author.ID)
			if len(currentReminders) >= 25 {
				return "You can have a maximum of 25 active reminders, list your reminders with the `reminders` command", nil
			}

			fromNow := parsed.Args[0].Value.(time.Duration)

			durString := common.HumanizeDuration(common.DurationPrecisionSeconds, fromNow)
			when := time.Now().Add(fromNow)
			tStr := when.UTC().Format(time.RFC822)

			if when.After(time.Now().Add(time.Hour * 24 * 366)) {
				return "Can be max 365 days from now...", nil
			}

			_, err := NewReminder(parsed.Msg.Author.ID, parsed.GS.ID, parsed.CS.ID, parsed.Args[1].Str(), when)
			if err != nil {
				return nil, err
			}

			return "Set a reminder in " + durString + " from now (" + tStr + ")\nView reminders with the reminders command", nil
		},
	},
	{
		CmdCategory: commands.CategoryTool,
		Name:        "Reminders",
		Description: "Lists your active reminders",
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			currentReminders, err := GetUserReminders(parsed.Msg.Author.ID)
			if err != nil {
				return nil, err
			}

			return "Your reminders:\n" + stringReminders(currentReminders, false) + "\nRemove a reminder with `delreminder/rmreminder (id)` where id is the first number for each reminder above", nil
		},
	},
	{
		CmdCategory:         commands.CategoryTool,
		Name:                "CReminders",
		Description:         "Lists reminders in channel, only users with 'manage server' permissions can use this.",
		RequireDiscordPerms: []int64{discordgo.PermissionManageChannels},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			currentReminders, err := GetChannelReminders(parsed.CS.ID)
			if err != nil {
				return nil, err
			}

			return "Reminders in this channel:\n" + stringReminders(currentReminders, true) + "\nRemove a reminder with `delreminder/rmreminder (id)` where id is the first number for each reminder above", nil
		},
	},
	{
		CmdCategory:  commands.CategoryTool,
		Name:         "DelReminder",
		Aliases:      []string{"rmreminder"},
		Description:  "Deletes a reminder. You can delete reminders from other users provided you are running this command in the same guild the reminder was created in and have the Manage Channel permission in the channel the reminder was created in.",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			{Name: "ID", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			var reminder Reminder
			err := common.GORM.Where(parsed.Args[0].Int()).First(&reminder).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return "No reminder by that id found", nil
				}

				return "Error retrieving reminder", err
			}

			// Check perms
			if reminder.UserID != discordgo.StrID(parsed.Msg.Author.ID) {
				if reminder.GuildID != parsed.GS.ID {
					return "You can only delete reminders that are not your own in the guild the reminder was originally created", nil
				}

				ok, err := bot.AdminOrPermMS(reminder.ChannelIDInt(), parsed.MS, discordgo.PermissionManageChannels)
				if err != nil {
					return nil, err
				}
				if !ok {
					return "You need manage channel permission in the channel the reminder is in to delete reminders that are not your own", nil
				}
			}

			// Do the actual deletion
			err = common.GORM.Delete(reminder).Error
			if err != nil {
				return nil, err
			}

			// Check if we should remove the scheduled event
			currentReminders, err := GetUserReminders(reminder.UserIDInt())
			if err != nil {
				return nil, err
			}

			delMsg := fmt.Sprintf("Deleted reminder **#%d**: '%s'", reminder.ID, limitString(reminder.Message))

			// If there is another reminder with the same timestamp, do not remove the scheduled event
			for _, v := range currentReminders {
				if v.When == reminder.When {
					return delMsg, nil
				}
			}

			return delMsg, nil
		},
	},
}

func stringReminders(reminders []*Reminder, displayUsernames bool) string {
	var out strings.Builder
	for _, v := range reminders {
		parsedCID, _ := strconv.ParseInt(v.ChannelID, 10, 64)

		t := time.Unix(v.When, 0)
		timeFromNow := common.HumanizeTime(common.DurationPrecisionMinutes, t)
		tStr := t.Format(time.RFC822)
		if !displayUsernames {
			channel := "<#" + discordgo.StrID(parsedCID) + ">"
			out.WriteString("**" + strconv.FormatUint(uint64(v.ID), 10) + "**: " + channel + ": '" + limitString(v.Message) + "' - " + timeFromNow + " from now (" + tStr + ")")
		} else {
			member, _ := bot.GetMember(v.GuildID, v.UserIDInt())
			username := "Unknown user"
			if member != nil {
				username = member.Username
			}
			out.WriteString("**" + strconv.FormatUint(uint64(v.ID), 10) + "**: " + username + ": '" + limitString(v.Message) + "' - " + timeFromNow + " from now (" + tStr + ")")
		}
	}

	return out.String()
}

func checkUserScheduledEvent(evt *seventsmodels.ScheduledEvent, data interface{}) (retry bool, err error) {
	// !important! the evt.GuildID can be 1 in cases where it was migrated from the legacy scheduled event system
	userID := *data.(*int64)

	reminders, err := GetUserReminders(userID)
	if err != nil {
		return true, err
	}

	now := time.Now()
	nowUnix := now.Unix()
	for _, v := range reminders {
		if v.When <= nowUnix {
			err := v.Trigger()
			if err != nil {
				// possibly try again
				return scheduledevents2.CheckDiscordErrRetry(err), err
			}
		}
	}

	return false, nil
}

func migrateLegacyScheduledEvents(t time.Time, data string) error {
	split := strings.Split(data, ":")
	if len(split) < 2 {
		logger.Error("invalid check user scheduled event: ", data)
		return nil
	}

	parsed, _ := strconv.ParseInt(split[1], 10, 64)

	return scheduledevents2.ScheduleEvent("reminders_check_user", 1, t, parsed)
}

func limitString(s string) string {
	if utf8.RuneCountInString(s) < 50 {
		return s
	}

	runes := []rune(s)
	return string(runes[:47]) + "..."
}
