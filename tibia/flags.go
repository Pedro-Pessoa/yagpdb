package tibia

import (
	"fmt"
	"sync"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/premium"
)

type TibiaFlags struct {
	common.SmallModel

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

func (m *TibiaFlags) TableName() string {
	return "tibia_flags"
}

func SetServerDeathChannel(server int64, channel int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	flags.SendDeaths = true
	flags.ChannelDeaths = channel

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao definir o canal.")
	}

	return fmt.Sprintf("Tudo certo! As notificações de mortes serão enviadas neste canal agora! <#%d>", channel), nil
}

func SetServerUpdatesChannel(server int64, channel int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	flags.SendUpdates = true
	flags.ChannelUpdates = channel

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao definir o canal.")
	}

	return fmt.Sprintf("Tudo certo! As notificações de players serão enviadas neste canal agora! <#%d>", channel), nil
}

func ToggleDeaths(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if !alreadySet {
		return "O mundo deste servidor ainda não foi decidido.", nil
	}

	if flags.SendDeaths {
		flags.SendDeaths = false
		err = common.GORM.Save(&flags).Error
		if err != nil {
			return "", err
		}

		return "Tudo certo! Não irei mais enviar notificações de mortes de players.", nil
	}

	flags.SendDeaths = true
	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", err
	}

	return "Beleza! Se algum dos players acompanhados morrerem eu irei avisar!", nil
}

func ToggleUpdates(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if !alreadySet {
		return "O mundo deste servidor ainda não foi decidido.", nil
	}

	if flags.SendUpdates {
		flags.SendUpdates = false
		err = common.GORM.Save(&flags).Error
		if err != nil {
			return "", err
		}

		return "Tudo certo! Não irei mais enviar notificações de players.", nil
	}

	flags.SendUpdates = true
	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", err
	}

	return "Beleza! Agora irei enviar notícias dos players acompanhados!", nil
}

func GetServerWorld(server int64, nameOnly bool) (string, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{WorldIsSet: true, ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if alreadySet {
		if !nameOnly {
			return fmt.Sprintf("O mundo deste servidor é **%s**!", flags.World), nil
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
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if alreadySet {
		return fmt.Sprintf("A guild deste servidor é **%s**!", flags.Guild), nil
	}

	return "A guild deste servidor ainda não foi definida!", nil
}

func SetServerWorld(world string, server int64, isAdmin bool) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if flags.WorldIsSet && !isAdmin {
		return "O mundo deste servidor já foi definido.", nil
	}

	mundo, err := GetWorld(world)
	if err != nil {
		return "", err
	}

	if len(mundo.World.WorldInformation.CreationDate) == 0 {
		return "Esse mundo não existe.", nil
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	hasTracks := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if hasTracks {
		common.GORM.Delete(&tracking)
	}

	flags.World = mundo.World.WorldInformation.Name
	flags.ServerID = server
	flags.WorldIsSet = true

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao definir o mundo.")
	}

	return fmt.Sprintf("Tudo certo! O mundo deste servidor agora é: **%s**", mundo.World.WorldInformation.Name), nil
}

var (
	flagwg sync.WaitGroup
	pool   chan struct{}
)

func SetServerGuild(guild string, server int64, isAdmin bool, memberCount int) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
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
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
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
		return "", errors.WithMessage(err, "Erro ao desserializar os tracks")
	}

	hunteds, err := deserializeValue(tracking.Hunteds)
	if err != nil {
		return "", errors.WithMessage(err, "Erro ao desserializar os hunteds")
	}

	isPremium, _ := premium.IsGuildPremium(server)
	limit := getServerLimit(memberCount, isPremium)
	if (len(deserialized) + len(hunteds) + len(tracks)) >= limit {
		return "Você não pode definir uma guild para este servidor por que o limite de chars que podem ser acompanhados já foi atingido.", nil
	}
	loopCap := limit - len(deserialized) - len(hunteds) - len(tracks)

	serverWorld, err := GetServerWorld(server, true)
	if err != nil {
		return "", err
	}

	cla, err := GetTibiaSpecificGuild(guild)
	if err != nil {
		return "", err
	}

	if serverWorld != "" {
		if cla.World != serverWorld {
			return fmt.Sprintf("Essa guild não pode ser definida como a guild do servidor porque ela não é do mundo **%s**!", serverWorld), nil
		}
	} else {
		return "O mundo deste servidor ainda não foi definido", nil
	}

	if len(cla.Members) == 0 {
		return "Essa guild não tem membros e por isso não pode ser acompanhada", nil
	}

	fila := make(chan InternalChar, len(cla.Members))
	pool = make(chan struct{}, 120)
	counterOut := 0
	for _, k := range cla.Members {
		flagwg.Add(1)
		go charFromMember(k.Name, fila, deserialized)
		counterOut += 1
	}

	flagwg.Wait()
	close(fila)
	close(pool)
	logger.Infof("Counter: %d", counterOut)

	counter := 0
	broke := false
	for e := range fila {
		if counter < loopCap {
			deserialized = append(deserialized, e)
			counter += 1
		} else {
			broke = true
			break
		}
	}

	goback, err := serializeValue(deserialized)
	if err != nil {
		return "", err
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
		return fmt.Sprintf("Tudo certo! A guild deste servidor agora é: **%s**\nDevido ao limite de chars que podem ser acompanhados neste servidor (%d), apenas %d membros da guild de %d estã sendo acompanhados!", cla.Name, limit, counter, len(cla.Members)), nil
	}

	return fmt.Sprintf("Tudo certo! A guild deste servidor agora é **%s** e todos os chars dela estão sendo acompanhados!", cla.Name), nil
}

func charFromMember(member string, channel chan InternalChar, tracks []InternalChar) {
	defer flagCleanUp()
	pool <- struct{}{}
	defer func() { <-pool }()
	char, err := GetTibiaChar(member, false)
	if err != nil {
		logger.Errorf("error 1 happened: %#v", err)
		return
	}
	for _, k := range tracks {
		if k.Name == char.Name {
			return
		}
	}
	channel <- *char
	return
}

func flagCleanUp() {
	if r := recover(); r != nil {
		logger.Infof("Recovered at: %v", r)
	}
	defer flagwg.Done()
}

func DeleteServerWorld(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	common.GORM.Delete(&tracking)

	flags.World = ""
	flags.WorldIsSet = false
	flags.ServerID = server

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao apagar o mundo deste server.")
	}

	return fmt.Sprintf("Tudo certo! O mundo do server **%d** foi removido!", server), nil
}

func DeleteServerGuild(server int64) (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&TibiaFlags{ServerID: server}).First(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	flags.Guild = ""
	flags.GuildIsSet = false
	flags.ServerID = server

	err = common.GORM.Save(&flags).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao apagar a guild deste server.")
	}

	return fmt.Sprintf("Tudo certo! A guild do server **%d** foi removida!", server), nil
}

func DeleteAll() (interface{}, error) {
	flags := TibiaFlags{}
	err := common.GORM.Where(&flags).Delete(&flags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	tracks := TibiaTracking{}
	err = common.GORM.Where(&tracks).Delete(&tracks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	table := ScanTable{}
	err = common.GORM.Where(&table).Delete(&table).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	return "Todas as databases foram apagadas!", nil
}
