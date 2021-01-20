package tibia

/*
import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate/v2"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/eventsystem"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/common/scheduledevents2"
	seventsmodels "github.com/jonas747/yagpdb/common/scheduledevents2/models"
)

type RespTable struct {
	sync.Mutex

	ServerID         int64 `gorm:"primary_key"`
	EnableRespSystem bool
	Channel          int64
	ChannelIsSet     bool
	IsPaused         bool
	Msgs             []int64
	DefaultDur       int
	RespBetterRole   int64
	RespBetterDur    int
	RespHighRole     int64
	RespHighDur      int
	IsRoleOneSet     bool
	IsRoleTwoSet     bool
	IsDefaultRoleSet bool
	DefaultRole      int64
	IsModRoleSet     bool
	ModRole          int64
}

type RespHandler struct {
	RespID int `json:"resp_id"`
}

type Respawns struct {
	sync.Mutex

	RespID      int `gorm:"primary_key"`
	IsRespHigh  bool
	Queue       []int64
	Start       *time.Time
	Pause       *time.Time
	IsPaused    bool
	IsRunning   bool
	CurrentUser *int64
}

func CreateRespawnMsg(cmdData *dcmd.Data, id int, guildID, userID int64, second bool) (string, error) {
	respTable := RespTable{}
	err := common.GORM.Where(&RespTable{ServerID: guildID}).First(&respTable).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", nil
	}

	if err == gorm.ErrRecordNotFound {
		return "O sistema de respawn não está ativo neste servidor.", nil
	}
	respTable.Lock()
	defer respTable.Unlock()

	if !respTable.ChannelIsSet || !respTable.EnableRespSystem || !respTable.IsDefaultRoleSet {
		return "O sistema de respawn não está ativo neste servidor.", nil
	}

	if respTable.IsPaused {
		return "O respawn está pausado.", nil
	}

	var member *dstate.MemberState
	var hasDRole, hasBRole, hasHRole, isMod, isOwner bool

	if cmdData.GS != nil {
		isOwner = userID == cmdData.GS.Guild.OwnerID
	}

	if cmdData != nil {
		member, _ = bot.GetMember(guildID, userID)
		if member == nil {
			logger.Errorf("Weird member error on guild %d", guildID)
			return "", nil
		}
		if !isOwner {
		OUTER:
			for _, r := range member.Roles {
				switch r {
				case respTable.DefaultRole:
					hasDRole = true
				case respTable.RespBetterRole:
					hasBRole = true
				case respTable.RespHighRole:
					hasHRole = true
				case respTable.ModRole:
					isMod = true
					break OUTER
				}
			}
		}
		if !hasDRole && !hasBRole && !hasHRole && !isMod && !isOwner {
			return "Você não tem permissão para usar o sistema de respawn", nil
		}
	}

	return "", nil
}

func HandleRespawns(evt *seventsmodels.ScheduledEvent, data interface{}) (retry bool, err error) {
	dataCast := data.(*Respawns)

	gs := bot.State.Guild(true, evt.GuildID)
	if gs == nil {
		return false, nil
	}

	resp := Respawns{}
	err = common.GORM.Where(&Respawns{RespID: dataCast.RespID}).First(&resp).Error
	if err != nil {
		return false, nil
	}

	if resp.IsPaused {
		time.Sleep(time.Minute * 2)
		return true, nil
	}

	gs.RLock()
	defer gs.RUnlock()
	resp.Lock()
	defer resp.Unlock()

	respTable := RespTable{}
	err = common.GORM.Where(&RespTable{ServerID: evt.GuildID}).First(&respTable).Error
	if err != nil {
		return false, nil
	}
	respTable.Lock()
	defer respTable.Unlock()

	if len(resp.Queue) == 0 {
		resp.Start = nil
		resp.IsRunning = false
		err = common.GORM.Save(&resp).Error
		if err != nil {
			logger.Errorf("Error 1 handling respawn: %#v", err)
		}
		return false, nil
	}

	resp.Start = func() *time.Time {
		a := time.Now()
		return &a
	}()
	resp.CurrentUser = &resp.Queue[0]
	resp.Queue = func(s []int64, i int) []int64 {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}(resp.Queue, 0)

	nextState := gs.MemberCopy(true, *resp.CurrentUser)
	if nextState == nil {
		err = scheduledevents2.ScheduleEvent("resp_handler", evt.GuildID, time.Now().Add(time.Minute), &RespHandler{
			RespID: dataCast.RespID,
		})

		if err != nil {
			logger.Errorf("Error 2 handling respawn: %#v", err)
		}

		err = common.GORM.Save(&respTable).Error
		if err != nil {
			logger.Errorf("Error 3 handling respawn: %#v", err)
		}

		return false, nil
	}

	var duration time.Time
	switch {
	case common.ContainsInt64Slice(nextState.Roles, respTable.RespHighRole), common.ContainsInt64Slice(nextState.Roles, respTable.ModRole):
		duration = time.Now().Add(time.Minute * time.Duration(respTable.RespHighDur))
	case common.ContainsInt64Slice(nextState.Roles, respTable.RespBetterRole):
		duration = time.Now().Add(time.Minute * time.Duration(respTable.RespBetterDur))
	default:
		duration = time.Now().Add(time.Minute * time.Duration(respTable.DefaultDur))
	}

	err = scheduledevents2.ScheduleEvent("resp_handler", evt.GuildID, duration, &RespHandler{
		RespID: dataCast.RespID,
	})

	if err != nil {
		logger.Errorf("Error 4 handling respawn: %#v", err)
	}

	err = common.GORM.Save(&respTable).Error
	if err != nil {
		logger.Errorf("Error 5 handling respawn: %#v", err)
	}

	return false, nil
}

func HandleMessageDelete(evt *eventsystem.EventData) {
	m := evt.MessageDelete()

	botUser := common.BotUser
	if botUser == nil || (botUser.ID != m.Author.ID) || len(m.Embeds) == 0 {
		return
	}

	resptable := RespTable{}
	err := common.GORM.Where(&RespTable{ServerID: evt.GS.ID}).First(&resptable).Error
	if err != nil {
		return
	}
	resptable.Lock()
	defer resptable.Unlock()

	if !resptable.ChannelIsSet || resptable.Channel != m.ChannelID {
		return
	}

	match := false
	index := 0
	for k, v := range resptable.Msgs {
		if v == m.ID {
			match = true
			index = k
			break
		}
	}

	if !match {
		return
	}

	newM, err := common.BotSession.ChannelMessageSend(resptable.Channel, m.Content)
	logger.Infof("Respawn msg resent on Guild %d", evt.GS.ID)
	if err != nil {
		return
	}

	newMsgs := func(s []int64, i int) []int64 {
		s[i] = s[len(s)-1]
		return s[:len(s)-1]
	}(resptable.Msgs, index)

	newMsgs = append(newMsgs, newM.ID)
	resptable.Msgs = newMsgs

	common.GORM.Save(&resptable)
}

var RespCommands = []*commands.YAGCommand{
	{
		CmdCategory: commands.CategoryTibia,
		Name:        "EnableResp",
		Aliases:     []string{"habilitarrespawn", "habilitarresp", "startresp", "hresp", "er"},
		Description: "Adiciona você à fila do respawn especificado.",
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			hasPerms, err := bot.AdminOrPermMS(data.CS.ID, data.MS, discordgo.PermissionKickMembers)
			if err != nil || !hasPerms {
				return "", nil
			}
			return "", nil
		},
	},
	{
		CmdCategory:  commands.CategoryTibia,
		Name:         "Resp",
		Description:  "Adiciona você à fila do respawn especificado.",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			{Name: "Respawn ID", Type: dcmd.Int},
		},
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			respID := data.Args[0].Int()
			found := false
			var respawnName string
			for k, v := range RespList {
				if k == respID {
					found = true
					respawnName = v
					break
				}
			}
			if !found {
				return "Número de respawn inválido.", nil
			}

			return respawnName, nil
		},
	},
}
*/
