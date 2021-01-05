package tibia

import (
	"bytes"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/yagpdb/common"
	"github.com/vmihailenco/msgpack"
)

type TibiaTracking struct {
	common.SmallModel

	ServerID int64 `gorm:"primary_key"`
	Tracks   []byte
	Guild    []byte
	Hunteds  []byte
}

func (m *TibiaTracking) TableName() string {
	return "tibia_tracking"
}

func TrackChar(char string, server int64, memberCount int, isPremium bool, isHunted bool) (interface{}, error) {
	getWorld, err := GetServerWorld(server, true)
	if err != nil {
		return "", err
	}

	if getWorld == "" {
		return "O mundo deste servidor ainda não foi definido!", nil
	}

	tracking := TibiaTracking{}
	err = common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var structure []byte
	var deserialized, check []InternalChar

	if !alreadySet {
		tracking.ServerID = server
		tracking.Tracks = structure
		tracking.Hunteds = structure
	}

	if isHunted {
		deserialized, err = deserializeValue(tracking.Hunteds)
		if err != nil {
			return "", err
		}

		check, err = deserializeValue(tracking.Tracks)
		if err != nil {
			return "", err
		}
	} else {
		deserialized, err = deserializeValue(tracking.Tracks)
		if err != nil {
			return "", err
		}

		check, err = deserializeValue(tracking.Hunteds)
		if err != nil {
			return "", err
		}
	}

	guild, err := deserializeValue(tracking.Guild)
	if err != nil {
		return "", err
	}

	if (len(deserialized) + len(check) + len(guild)) >= getServerLimit(memberCount, isPremium) {
		return "Infelizmente este servidor já chegou no limite de chares que podem ser acompanhados", nil
	}

	insideChar, err := GetTibiaChar(char, true)
	if err != nil {
		return "", err
	}

	if getWorld != insideChar.World {
		return fmt.Sprintf("Você só pode fazer track de chars do mundo **%s**", getWorld), nil
	}

	already := "Esse char já está sendo acompanhado!"

	if len(deserialized) > 0 {
		for _, v := range deserialized {
			if strings.ToLower(v.Name) == strings.ToLower(insideChar.Name) {
				return already, nil
			}
		}
	}

	if len(check) > 0 {
		for _, e := range check {
			if strings.ToLower(e.Name) == strings.ToLower(insideChar.Name) {
				return already, nil
			}
		}
	}

	if len(guild) > 0 {
		for _, k := range guild {
			if strings.ToLower(k.Name) == strings.ToLower(insideChar.Name) {
				return already, nil
			}
		}
	}

	deserialized = append(deserialized, *insideChar)
	goback, err := serializeValue(deserialized)
	if err != nil {
		return "", err
	}

	if isHunted {
		tracking.Hunteds = goback
	} else {
		tracking.Tracks = goback
	}

	err = common.GORM.Save(&tracking).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao salvar este char.")
	}

	return fmt.Sprintf("Tudo certo! Agora o char **%s** está sendo acompanhado!", insideChar.Name), nil
}

func UnTrackChar(char string, server int64, hunted bool, guild bool) (interface{}, error) {
	tracking := TibiaTracking{}
	err := common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if !alreadySet {
		return "Nenhum char está sendo acompanhado neste servidor.", nil
	}

	var deserialized []InternalChar

	if guild {
		deserialized, err = deserializeValue(tracking.Guild)
		if err != nil {
			return "", err
		}
	} else if hunted {
		deserialized, err = deserializeValue(tracking.Hunteds)
		if err != nil {
			return "", err
		}
	} else {
		deserialized, err = deserializeValue(tracking.Tracks)
		if err != nil {
			return "", err
		}
	}

	if len(deserialized) == 0 {
		return "Esse char não está sendo acompanhado.", nil
	}

	found := false
	var index int
	for k, v := range deserialized {
		if strings.ToLower(v.Name) == strings.ToLower(char) {
			found = true
			index = k
			break
		}
	}

	if !found {
		return "Esse char não está sendo acompanhado!", nil
	}

	deserialized = removeFromSlice(deserialized, index)

	goback, err := serializeValue(deserialized)
	if err != nil {
		return "", err
	}

	if guild {
		tracking.Guild = goback
	} else if hunted {
		tracking.Hunteds = goback
	} else {
		tracking.Tracks = goback
	}

	err = common.GORM.Save(&tracking).Error
	if err != nil {
		return "", errors.WithMessage(err, "Algo deu errado ao apagar este char.")
	}

	return fmt.Sprintf("Tudo certo!! O char **%s** não está mais sendo acompanhado!", char), nil
}

func serializeValue(v []InternalChar) ([]byte, error) {
	var b bytes.Buffer
	enc := msgpack.NewEncoder(&b)
	err := enc.Encode(v)
	return b.Bytes(), err
}

func deserializeValue(v []byte) ([]InternalChar, error) {
	var out []InternalChar
	if len(v) == 0 {
		return out, nil
	}
	err := msgpack.Unmarshal(v, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

func removeFromSlice(s []InternalChar, i int) []InternalChar {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getServerLimit(memberCount int, isPremium bool) int {
	value := 100
	if isPremium {
		value = 6000
	} else {
		if memberCount < 31 {
			value = 50
		} else if memberCount > 149 {
			value = 350
		}
	}
	return value
}

func GetTracks(server int64) ([]InternalChar, error) {
	tracking := TibiaTracking{}
	err := common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ErrWithCaller(err)
	}

	if !alreadySet || (len(tracking.Tracks) == 0 && len(tracking.Guild) == 0) {
		return nil, nil
	}

	deserialized, err := deserializeValue(tracking.Tracks)
	if err != nil {
		return nil, err
	}

	guild, err := deserializeValue(tracking.Guild)
	if err != nil {
		return nil, err
	}

	deserialized = append(deserialized, guild...)

	return deserialized, nil
}

func GetHuntedList(server int64) ([]InternalChar, error) {
	tracking := TibiaTracking{}
	err := common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, common.ErrWithCaller(err)
	}

	if !alreadySet || len(tracking.Hunteds) == 0 {
		return nil, nil
	}

	deserialized, err := deserializeValue(tracking.Hunteds)
	if err != nil {
		return nil, err
	}

	return deserialized, nil
}

func FindAll() ([]TibiaTracking, error) {
	flags := []TibiaTracking{}
	err := common.GORM.Find(&flags).Error
	if err != nil {
		return nil, common.ErrWithCaller(err)
	}

	return flags, nil
}

func FindAllGuilds() ([]TibiaFlags, error) {
	guilds := []TibiaFlags{}
	err := common.GORM.Find(&guilds).Error
	if err != nil {
		return nil, err
	}

	return guilds, nil
}

func DeleteTracks(server int64, hunted bool, guild bool, all bool) (string, error) {
	tracking := TibiaTracking{}
	err := common.GORM.Where(&TibiaTracking{ServerID: server}).First(&tracking).Error
	alreadySet := err != gorm.ErrRecordNotFound
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", common.ErrWithCaller(err)
	}

	if !alreadySet || (len(tracking.Guild) == 0 && len(tracking.Hunteds) == 0 && len(tracking.Tracks) == 0) {
		return "Nenhum char está sendo trackeado nesse servidor.", nil
	}

	if all {
		tracking.Guild = []byte{}
		tracking.Hunteds = []byte{}
		tracking.Tracks = []byte{}
	} else if guild {
		tracking.Guild = []byte{}
	} else if hunted {
		tracking.Hunteds = []byte{}
	} else {
		tracking.Tracks = []byte{}
	}

	err = common.GORM.Save(&tracking).Error
	if err != nil {
		return "", err
	}

	return "Tudo certo!", nil
}
