package tibia

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/premium"
)

type ScanTable struct {
	common.SmallModel

	RunScan bool
}

func (st *ScanTable) TableName() string {
	return "scan_table"
}

type DataStore struct {
	sync.Mutex
	cache   []InternalChar
	counter int
	total   int
}

func New() *DataStore {
	return &DataStore{
		counter: 0,
		total:   0,
		cache:   []InternalChar{},
	}
}

func (ds *DataStore) set(value InternalChar) {
	ds.Lock()
	defer ds.Unlock()
	ds.cache = append(ds.cache, value)
}

func (ds *DataStore) get(name string) *InternalChar {
	ds.Lock()
	defer ds.Unlock()
	for _, e := range ds.cache {
		if e.Name == name {
			return &e
		}
	}
	return nil
}

func (ds *DataStore) flush() {
	ds.cache = []InternalChar{}
	ds.counter = 0
	ds.total = 0
}

var (
	masterwg, trackwg, msgswg, msgshuntedwg, guildswg, updatewg, guildwg sync.WaitGroup
	trackpool, updatepool                                                chan struct{}
)

func (ds *DataStore) scanTracks() {
	defer masterwg.Done()
	start := time.Now()
	trackpool = make(chan struct{}, 500)
	logger.Infof("Tracking starting... %v", start)

	guilds, err := FindAllGuilds()
	if err != nil || len(guilds) == 0 {
		logger.Info("As guilds não foram encontradas!")
	} else {
		guildTimer := time.Now()
		logger.Info("Guilds encontradas, atualizando...")
		for _, g := range guilds {
			guildswg.Add(1)
			go updateGuild(g, ds)
		}

		guildswg.Wait()
		logger.Infof("Guilds atualizadas em %vs... Continuando...", time.Since(guildTimer).Seconds())
	}

	tracks, err := FindAll()
	if err != nil || len(tracks) == 0 {
		logger.Infof("Nothing to scan. Scan concluido em %vs", time.Since(start).Seconds())
		return
	}

	logger.Info("Found all, tracking...")

	for _, v := range tracks {
		trackwg.Add(1)
		go ds.trackRoutine(v)
	}

	trackwg.Wait()
	logger.Infof("Scan concluido em %vs!\nHTTP Requests made: %d\nTotal size: %d", time.Since(start).Seconds(), ds.counter, ds.total)
	ds.flush()
}

func (ds *DataStore) trackingController() {
	table := ScanTable{}
	done := make(chan bool, 1)
	go func() {
		for {
			err := common.GORM.Where(&ScanTable{}).First(&table).Error
			alreadySet := err != gorm.ErrRecordNotFound
			if err != nil && alreadySet {
				logger.Errorf("Err on trackingController: %v", err)
				return
			}

			if !alreadySet || !table.RunScan {
				done <- true
			}

			select {
			case <-done:
				close(done)
				logger.Info("Tracking done")
				return
			default:
				masterwg.Add(1)
				logger.Info("Running Track")
				ds.scanTracks()
			}

			masterwg.Wait()

			time.Sleep(10 * time.Second)
		}
	}()
}

func StartLoop() (string, error) {
	table := ScanTable{}
	err := common.GORM.Where(&table).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if alreadySet && table.RunScan {
		return "O tracking já está rolando.", nil
	}

	table.RunScan = true

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	store := New()
	store.trackingController()

	return "Tudo certo! O tracking está rolando", nil
}

func StopLoop() (string, error) {
	table := ScanTable{}
	err := common.GORM.Where(&table).First(&table).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if !alreadySet || !table.RunScan {
		return "O tracking não está rolando ainda.", nil
	}

	table.RunScan = false

	err = common.GORM.Save(&table).Error
	if err != nil {
		return "", err
	}

	return "Tudo certo! O tracking foi pausado.", nil
}

func (ds *DataStore) trackRoutine(input TibiaTracking) {
	defer trackCleanUp()
	trackpool <- struct{}{}
	defer func() { <-trackpool }()

	g := bot.State.Guild(true, input.ServerID)
	if len(input.Tracks) == 0 || g == nil {
		return
	}

	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: input.ServerID}).First(&flags).Error
	if err != nil {
		return
	}

	changed := false

	if len(input.Tracks) > 0 {
		changed = true
		deserialized, err := deserializeValue(input.Tracks)
		if err != nil {
			return
		}

		channel := make(chan InternalChar, len(deserialized))

		for k, v := range deserialized {
			msgswg.Add(1)
			go ds.msgsRoutine(v, k, channel, flags, false, false)
		}

		msgswg.Wait()
		close(channel)

		var output []InternalChar
		for e := range channel {
			output = append(output, e)
		}

		serialized, err := serializeValue(output)
		if err != nil {
			return
		}

		input.Tracks = serialized
	}

	if len(input.Hunteds) > 0 {
		changed = true
		deserialized, err := deserializeValue(input.Hunteds)
		if err != nil {
			return
		}

		channel := make(chan InternalChar, len(deserialized))

		for k, v := range deserialized {
			msgshuntedwg.Add(1)
			go ds.msgsRoutine(v, k, channel, flags, true, false)
		}

		msgshuntedwg.Wait()
		close(channel)

		var output []InternalChar
		for e := range channel {
			output = append(output, e)
		}

		serialized, err := serializeValue(output)
		if err != nil {
			return
		}

		input.Hunteds = serialized
	}

	if len(input.Guild) > 0 {
		changed = true
		deserialized, err := deserializeValue(input.Guild)
		if err != nil {
			return
		}

		channel := make(chan InternalChar, len(deserialized))

		for k, v := range deserialized {
			guildwg.Add(1)
			go ds.msgsRoutine(v, k, channel, flags, false, true)
		}

		guildwg.Wait()
		close(channel)

		var output []InternalChar
		for e := range channel {
			output = append(output, e)
		}

		serialized, err := serializeValue(output)
		if err != nil {
			return
		}

		input.Guild = serialized
	}

	if !changed {
		return
	}

	common.GORM.Save(&input)
}

func (ds *DataStore) msgsRoutine(input InternalChar, k int, channel chan InternalChar, flags TibiaFlags, areHunteds bool, isGuild bool) {
	defer msgsCleanUp(areHunteds, isGuild)
	trackpool <- struct{}{}
	defer func() { <-trackpool }()
	var output, deathsoutput string
	var char InternalChar
	found := false
	ds.total += 1
	currentChar := ds.get(input.Name)

	if currentChar != nil {
		char = *currentChar
		found = true
	}

	if !found {
		income, err := GetTibiaChar(input.Name, true)
		ds.counter += 1
		if err != nil || income == nil {
			return
		}
		char = *income
		ds.set(char)
	}

	if char.Name != input.Name {
		output = fmt.Sprintf("%s\nParece que ele tava insatisfeito com o nome e agora se chama **%s**!!", output, char.Name)
		input.Name = char.Name
	}

	if char.Level != input.Level {
		if char.Level > input.Level {
			output = fmt.Sprintf("%s\nUPOOUU! Agora tá no level: **%d**!", output, char.Level)
		}

		input.Level = char.Level
	}

	if char.World != input.World {
		output = fmt.Sprintf("%s\nDesertor ou auxiliar de guerra? Ele fez uma viagem longa e atracou em **%s**!!", output, char.World)
		input.World = char.World
	}

	if char.Residence != input.Residence {
		output = fmt.Sprintf("%s\nNão tava gostando da cidade natal né? O que você tá achando de **%s**?", output, char.Residence)
		input.Residence = char.Residence
	}

	if char.AchievementPoints != input.AchievementPoints {
		if char.AchievementPoints > input.AchievementPoints {
			output = fmt.Sprintf("%s\nOlha o char lover ai!! Upando achievement, mano!! Agora tá com **%d** pontos!", output, char.AchievementPoints)
		}

		input.AchievementPoints = char.AchievementPoints
	}

	if char.Sex != input.Sex {
		output = fmt.Sprintf("%s\nMomento de inclusão! Esse char agora é **%s**!", output, char.Sex)
		input.Sex = char.Sex
	}

	if char.Married != input.Married {
		if char.Married != "Ninguém" {
			output = fmt.Sprintf("%s\nSe casou com **%s**!!", output, char.Married)
		} else {
			output = fmt.Sprintf("%s\nSe divorciou!!!", output)
		}

		input.Married = char.Married
	}

	if char.Guild != input.Guild {
		if char.Guild != "Sem guild" {
			output = fmt.Sprintf("%s\nAbre o olho ae em!! Mudou de guild e agora está na **%s**", output, char.Guild)
		} else {
			output = fmt.Sprintf("%s\nNão faz mais parte de guild nenhuma!", output)
		}

		input.Guild = char.Guild
	}

	if char.Rank != input.Rank {
		if char.Rank != "Sem guild" {
			output = fmt.Sprintf("%s\nMudou de cargo na guild e agora é um **%s**!!", output, char.Rank)
		}

		input.Rank = char.Rank
	}

	if !reflect.DeepEqual(char.Deaths, input.Deaths) {
		if len(char.Deaths) > 0 {
			if len(input.Deaths) > 0 {
				if char.Deaths[0] != input.Deaths[0] {
					deathsoutput = fmt.Sprintf("Level: **%d**\nMotivo: **%s**\nData: **%s**", char.Deaths[0].Level, char.Deaths[0].Reason, char.Deaths[0].Date)
				}
			} else {
				deathsoutput = fmt.Sprintf("Level: **%d**\nMotivo: **%s**\nData: **%s**", char.Deaths[0].Level, char.Deaths[0].Reason, char.Deaths[0].Date)
			}
		}

		input.Deaths = char.Deaths
	}

	channel <- input

	msgSend := &discordgo.MessageSend{
		AllowedMentions: discordgo.AllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers, discordgo.AllowedMentionTypeRoles, discordgo.AllowedMentionTypeEveryone},
		},
	}

	var title string
	if flags.SendUpdates && len(output) > 0 {
		if areHunteds {
			title = fmt.Sprintf("Tem novidade sobre %s (HUNTED)", input.Name)
		} else {
			title = fmt.Sprintf("Tem novidade sobre %s (FRIEND)", input.Name)
		}

		msgSend.Embed = &discordgo.MessageEmbed{
			Title:       title,
			Description: output,
			Color:       int(rand.Int63n(16777215)),
		}

		_, _ = common.BotSession.ChannelMessageSendComplex(flags.ChannelUpdates, msgSend)
	}

	if flags.SendDeaths && len(deathsoutput) > 0 {
		if areHunteds {
			title = fmt.Sprintf("ENEMY MORTO: %s", input.Name)
		} else {
			title = fmt.Sprintf("FRIEND MORTO: %s", input.Name)
		}

		msgSend.Embed = &discordgo.MessageEmbed{
			Title:       title,
			Description: deathsoutput,
			Color:       int(rand.Int63n(16777215)),
		}

		_, _ = common.BotSession.ChannelMessageSendComplex(flags.ChannelDeaths, msgSend)
	}
}

func updateGuild(g TibiaFlags, ds *DataStore) {
	defer guildswg.Done()
	trackpool <- struct{}{}
	defer func() { <-trackpool }()

	tracking := TibiaTracking{}
	err := common.GORM.Where(&TibiaTracking{ServerID: g.ServerID}).First(&tracking).Error
	if err != nil {
		return
	}

	state := bot.State.Guild(true, g.ServerID)
	if state == nil || !g.GuildIsSet {
		return
	}

	deserialized, err := deserializeValue(tracking.Guild)
	if err != nil {
		return
	}

	if len(deserialized) == 0 {
		return
	}

	guild, err := GetTibiaSpecificGuild(g.Guild)
	if err != nil {
		return
	}

	if len(guild.Members) == 0 {
		return
	}

	tracks, err := deserializeValue(tracking.Tracks)
	if err != nil {
		return
	}

	hunteds, err := deserializeValue(tracking.Hunteds)
	if err != nil {
		return
	}

	isPremium, _ := premium.IsGuildPremium(g.ServerID)
	limit := getServerLimit(state.Guild.MemberCount, isPremium)
	length := len(deserialized) + len(hunteds) + len(tracks)
	if length >= limit {
		return
	}
	loopCap := limit - length

	fila := make(chan InternalChar, len(guild.Members))
	updatepool = make(chan struct{}, 100)

	for _, k := range guild.Members {
		ds.total += 1
		if !alreadyTracked(deserialized, k) {
			updatewg.Add(1)
			go func() {
				defer updatewg.Done()
				updatepool <- struct{}{}
				defer func() { <-updatepool }()
				var char *InternalChar
				var err error
				if a := ds.get(k.Name); a == nil {
					char, err = GetTibiaChar(k.Name, true)
					ds.counter += 1
					if err != nil || char == nil {
						logger.Errorf("Error on update: %#v", err)
						return
					}
				} else {
					char = a
				}
				fila <- *char
			}()
		}
	}

	updatewg.Wait()

	counter := 0
	if len(fila) > 0 {
		for e := range fila {
			if counter < loopCap {
				deserialized = append(deserialized, e)
				counter++
			}

			if ds.get(e.Name) == nil {
				ds.set(e)
			}
		}
	}

	if counter >= loopCap {
		logger.Infof("Server chegou no limite: %d", g.ServerID)
	}

	goback, err := serializeValue(deserialized)
	if err != nil {
		logger.Error(err)
		return
	}

	tracking.Guild = goback
	common.GORM.Save(&tracking)
}

func alreadyTracked(list []InternalChar, e GuildMember) bool {
	for _, k := range list {
		if k.Name == e.Name {
			return true
		}
	}

	return false
}

func msgsCleanUp(a, b bool) {
	if r := recover(); r != nil {
		logger.Infof("Recovered at: %v", r)
	}

	switch {
	case a:
		defer msgshuntedwg.Done()
	case b:
		defer guildwg.Done()
	default:
		defer msgswg.Done()
	}
}

func trackCleanUp() {
	if r := recover(); r != nil {
		logger.Infof("Recovered at: %v", r)
	}

	defer trackwg.Done()
}
