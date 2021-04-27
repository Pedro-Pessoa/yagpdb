package tibia

import (
	"strconv"
	"sync"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/mediocregopher/radix/v3"

	"github.com/Pedro-Pessoa/tidbot/common"
	"github.com/Pedro-Pessoa/tidbot/premium"
)

type TibiaFlags struct {
	ServerID       int64 `gorm:"primary_key"`
	World          string
	Guild          string
	WorldIsSet     bool
	GuildIsSet     bool
	SendDeaths     bool
	ChannelDeaths  int64
	SendUpdates    bool
	ChannelUpdates int64
}

func (tf *TibiaFlags) TableName() string {
	return "tibia_flags"
}

func SetServerDeathChannel(server int64, channel int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "SetServerDeathChannel Error 1")
	}

	flags.ServerID = server
	flags.SendDeaths = true
	flags.ChannelDeaths = channel

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "SetServerDeathChannel Error 2")
	}

	return "Tudo certo! As notificações de mortes serão enviadas neste canal agora! <#" + strconv.FormatInt(channel, 10) + ">", nil
}

func SetServerUpdatesChannel(server int64, channel int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "SetServerUpdateChannel Error 1")
	}

	flags.ServerID = server
	flags.SendUpdates = true
	flags.ChannelUpdates = channel

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "SetServerUpdateChannel Error 2")
	}

	return "Tudo certo! As notificações de players serão enviadas neste canal agora! <#" + strconv.FormatInt(channel, 10) + ">", nil
}

func ToggleDeaths(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", errors.WithMessage(err, "ToggleDeaths Error 1")
	}

	if !alreadySet {
		return "O mundo deste servidor ainda não foi decidido.", nil
	}

	if flags.SendDeaths {
		flags.SendDeaths = false
		err = common.GORM.Save(&flags).Error
		if err != nil {
			return "", errors.WithMessage(err, "ToggleDeaths Error 2")
		}

		return "Tudo certo! Não irei mais enviar notificações de mortes de players.", nil
	}

	flags.SendDeaths = true
	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "ToggleDeaths Error 3")
	}

	return "Beleza! Se algum dos players acompanhados morrerem eu irei avisar!", nil
}

func ToggleUpdates(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", errors.WithMessage(err, "ToggleUpdates Error 1")
	}

	if !alreadySet {
		return "O mundo deste servidor ainda não foi decidido.", nil
	}

	if flags.SendUpdates {
		flags.SendUpdates = false
		err = common.GORM.Save(&flags).Error
		if err != nil {
			return "", errors.WithMessage(err, "ToggleUpdates Error 2")
		}

		return "Tudo certo! Não irei mais enviar notificações de players.", nil
	}

	flags.SendUpdates = true
	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "ToggleUpdates Error 3")
	}

	return "Beleza! Agora irei enviar notícias dos players acompanhados!", nil
}

func GetServerWorld(server int64, nameOnly bool) (string, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{WorldIsSet: true, ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", err
	}

	if alreadySet {
		if !nameOnly {
			return "O mundo deste servidor é **" + flags.World + "**!", nil
		}

		return flags.World, nil
	}

	if nameOnly {
		return "", nil
	}

	return "O mundo deste servidor ainda não foi definido!", nil
}

func GetServerGuild(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{GuildIsSet: true, ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", err
	}

	if alreadySet {
		return "A guild deste servidor é **" + flags.Guild + "**!", nil
	}

	return "A guild deste servidor ainda não foi definida!", nil
}

func SetServerWorld(world string, server int64, isAdmin bool) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "SetServerWorld Error 1")
	}

	if flags.WorldIsSet && !isAdmin {
		return "O mundo deste servidor já foi definido.", nil
	}

	mundo, err := getWorld(world)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerWorld Error 2")
	}

	if len(mundo.World.WorldInformation.CreationDate) == 0 {
		return "Esse mundo não existe.", nil
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	hasTracks := err != gorm.ErrRecordNotFound
	if err != nil && hasTracks {
		return "", errors.WithMessage(err, "SetServerWorld Error 3")
	}

	if hasTracks {
		err = common.GORM.Delete(&tracking).Error
		if err != nil {
			return "", errors.WithMessage(err, "SetServerWorld Error 4")
		}
	}

	flags.World = mundo.World.WorldInformation.Name
	flags.ServerID = server
	flags.WorldIsSet = true

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "SetServerWorld Error 5")
	}

	return "Tudo certo! O mundo deste servidor agora é: **" + mundo.World.WorldInformation.Name + "**", nil
}

var (
	flagwg sync.WaitGroup
	pool   chan struct{}
)

func SetServerGuild(guild string, server int64, isAdmin bool, memberCount int) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", errors.WithMessage(err, "SetServerGuild Error 1")
	}

	if !alreadySet {
		return "O mundo deste servidor ainda não foi definido.", nil
	}

	if flags.GuildIsSet && !isAdmin {
		return "A guild deste servidor já foi definida.", nil
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet = err != gorm.ErrRecordNotFound
	if err != nil && alreadySet {
		return "", errors.WithMessage(err, "SetServerGuild Error 2")
	}

	if !alreadySet {
		tracking.Tracks = []byte{}
		tracking.Hunteds = []byte{}
	}

	tracking.ServerID = server
	tracking.Guild = []byte{}

	deserialized, err := deserializeValue(tracking.Guild)
	if err != nil {
		return "", err
	}

	tracks, err := deserializeValue(tracking.Tracks)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerGuild Error 3")
	}

	hunteds, err := deserializeValue(tracking.Hunteds)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerGuild Error 4")
	}

	isPremium, _ := premium.IsGuildPremium(server)
	limit := getServerLimit(memberCount, isPremium)
	if (len(deserialized) + len(hunteds) + len(tracks)) >= limit {
		return "Você não pode definir uma guild para este servidor por que o limite de chars que podem ser acompanhados já foi atingido.", nil
	}
	loopCap := limit - len(deserialized) - len(hunteds) - len(tracks)

	serverWorld, err := GetServerWorld(server, true)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerGuild Error 5")
	}

	cla, err := GetTibiaSpecificGuild(guild)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerGuild Error 6")
	}

	if serverWorld != "" {
		if cla.World != serverWorld {
			return "Essa guild não pode ser definida como a guild do servidor porque ela não é do mundo **" + serverWorld + "**!", nil
		}
	} else {
		return "O mundo deste servidor ainda não foi definido", nil
	}

	if len(cla.Members) == 0 {
		return "Essa guild não tem membros e por isso não pode ser acompanhada", nil
	}

	fila := make(chan InternalChar, len(cla.Members))
	pool = make(chan struct{}, 120)
	var counterOut int

	for _, k := range cla.Members {
		flagwg.Add(1)
		go charFromMember(k.Name, fila, deserialized)
		counterOut++
	}

	flagwg.Wait()
	close(fila)
	close(pool)
	logger.Infof("Counter: %d", counterOut)

	var counter int
	var broke bool
	for e := range fila {
		if counter < loopCap {
			deserialized = append(deserialized, e)
			counter++
		} else {
			broke = true
			break
		}
	}

	goback, err := serializeValue(deserialized)
	if err != nil {
		return "", errors.WithMessage(err, "SetServerGuild Error 7")
	}

	tracking.Guild = goback
	err = common.GORM.Save(&tracking).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao salvar os chars da guild.")
	}

	flags.Guild = cla.Name
	flags.ServerID = server
	flags.GuildIsSet = true

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao definir a guild.")
	}

	if broke {
		return "Tudo certo! A guild deste servidor agora é: **" + cla.Name + "**\nDevido ao limite de chars que podem ser acompanhados neste servidor _**(" + strconv.Itoa(limit) + ")**_, apenas **" + strconv.Itoa(counter) + "** membros da guild de " + strconv.Itoa(len(cla.Members)) + " estã sendo acompanhados!", nil
	}

	return "Tudo certo! A guild deste servidor agora é **" + cla.Name + "** e todos os chars dela estão sendo acompanhados!", nil
}

func charFromMember(member string, channel chan InternalChar, tracks []InternalChar) {
	defer flagwg.Done()

	pool <- struct{}{}
	defer func() { <-pool }()

	char, err := GetTibiaChar(member, false)
	if err != nil || char == nil {
		logger.Errorf("error 1 on charFromMember happened: %#v || string err: %s", err, err)
		return
	}

	for _, k := range tracks {
		if k.Name == char.Name {
			return
		}
	}

	channel <- *char
}

func DeleteServerWorld(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "Delete Server World Error 1")
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "Delete Server World Error 2")
	}

	err = common.GORM.Delete(&tracking).Error
	if err != nil {
		return "", errors.WithMessage(err, "Delete Server World Error 3")
	}

	flags.World = ""
	flags.WorldIsSet = false
	flags.ServerID = server

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Delete Server World Error 4")
	}

	return "Tudo certo! O mundo do server **" + strconv.FormatInt(server, 10) + "** foi removido!", nil
}

func DeleteServerGuild(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteServerGuild Error 1")
	}

	flags.Guild = ""
	flags.GuildIsSet = false
	flags.ServerID = server

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao apagar a guild deste server.")
	}

	return "Tudo certo! A guild do server **" + strconv.FormatInt(server, 10) + "** foi removida!", nil
}

func DeleteAll() (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&flags).Delete(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteAll Error 1")
	}

	tracks := TibiaTracking{}
	err = common.GORM.Where(&tracks).Delete(&tracks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteAll Error 2")
	}

	table := ScanTable{}
	err = common.GORM.Where(&table).Delete(&table).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteAll Error 3")
	}

	news := NewsTable{}
	err = common.GORM.Where(&news).Delete(&news).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteAll Error 4")
	}

	inner := InnerNewsStruct{}
	err = common.GORM.Where(&inner).Delete(&inner).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", errors.WithMessage(err, "DeleteAll Error 5")
	}

	var i int
	err = common.RedisPool.Do(radix.Cmd(&i, "DEL", "news_guilds"))
	if err != nil {
		return "", errors.WithMessage(err, "DeleteAll Error 6")
	}

	return "Todas as databases foram apagadas!\n" + strconv.Itoa(i) + " entries no redis deletadas", nil
}