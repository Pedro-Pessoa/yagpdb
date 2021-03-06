package tibia

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/araddon/dateparse"
)

var (
	ErrInvalidName = errors.New("Você tem que especificar um char.")
	ErrSmallName   = errors.New("O nome fornecido é inválido.")
)

func GetTibiaChar(char string, update bool) (*InternalChar, error) {
	if len(char) <= 0 {
		return nil, ErrInvalidName
	} else if len(char) < 2 {
		return nil, ErrSmallName
	}

	tibia, err := GetChar(char)
	if err != nil {
		return nil, err
	}

	matched, _ := regexp.MatchString(`Character does not exist.`, tibia.Characters.Error)
	if matched {
		return nil, errors.New("Esse char não existe.")
	}

	level := tibia.Characters.Data.Level
	if update {
		world, err := GetWorld(tibia.Characters.Data.World)
		if err != nil {
			return nil, err
		}

		for _, v := range world.World.PlayersOnline {
			if v.Name == tibia.Characters.Data.Name {
				if v.Level > tibia.Characters.Data.Level {
					level = v.Level
				}
				break
			}
		}
	}

	comentario := "Char sem comentário"
	if len(tibia.Characters.Data.Comment) >= 1 {
		comentario = tibia.Characters.Data.Comment
	}

	lealdade := "Sem lealdade"
	if len(tibia.Characters.AccountInformation.LoyaltyTitle) > 0 {
		lealdade = tibia.Characters.AccountInformation.LoyaltyTitle
	}

	guild := "Sem guild"
	cargo := "Sem cargo"
	if len(tibia.Characters.Data.Guild.Name) >= 1 {
		guild = tibia.Characters.Data.Guild.Name
		cargo = tibia.Characters.Data.Guild.Rank
	}

	casado := "Ninguém"
	if len(tibia.Characters.Data.MarriedTo) >= 1 {
		casado = tibia.Characters.Data.MarriedTo
	}

	casa := "Nenhuma"
	if len(tibia.Characters.Data.House.Name) >= 1 {
		casa = tibia.Characters.Data.House.Name
	}

	criado := "Data escondida"
	if len(tibia.Characters.AccountInformation.Created.Date) > 0 {
		t, err := dateparse.ParseLocal(tibia.Characters.AccountInformation.Created.Date)
		if err != nil {
			return nil, err
		}
		criado = (t.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT")
	}

	mortes := []InternalDeaths{}
	if len(tibia.Characters.Deaths) > 0 {
		for _, v := range tibia.Characters.Deaths {
			t2, err := dateparse.ParseLocal(v.Date.Date)
			if err != nil {
				return nil, err
			}
			var insert InternalDeaths
			insert.Name = tibia.Characters.Data.Name
			insert.Level = v.Level
			insert.Reason = v.Reason
			insert.Date = (t2.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT")
			mortes = append(mortes, insert)
		}
	}

	output := InternalChar{}
	output.Name = tibia.Characters.Data.Name
	output.Level = level
	output.World = tibia.Characters.Data.World
	output.Vocation = tibia.Characters.Data.Vocation
	output.Residence = tibia.Characters.Data.Residence
	output.AccountStatus = tibia.Characters.Data.AccountStatus
	output.Status = strings.Title(tibia.Characters.Data.Status)
	output.Loyalty = lealdade
	output.AchievementPoints = tibia.Characters.Data.AchievementPoints
	output.Sex = strings.Title(tibia.Characters.Data.Sex)
	output.Married = casado
	output.Guild = guild
	output.Rank = cargo
	output.Comment = comentario
	output.CreatedAt = criado
	output.House = casa
	output.Deaths = mortes

	return &output, nil
}

func GetTibiaSpecificGuild(guildName string) (*InternalGuild, error) {
	if len(guildName) <= 0 {
		return nil, errors.New("Você tem que especificar uma guild.")
	} else if len(guildName) < 2 {
		return nil, ErrSmallName
	}

	guild, err := GetSpecificGuild(strings.Title(guildName))
	if err != nil {
		return nil, err
	}

	if len(guild.Guild.Error) >= 1 {
		return nil, errors.New("Essa guild não existe")
	}

	desc := "Guild sem descrição."
	if len(guild.Guild.Data.Description) >= 1 && len(guild.Guild.Data.Description) < 2048 {
		desc = guild.Guild.Data.Description
	}

	guildHall := "Nenhuma."
	if len(guild.Guild.Data.Guildhall.Name) > 1 {
		guildHall = fmt.Sprintf("**%s** que fica em %s", guild.Guild.Data.Guildhall.Name, guild.Guild.Data.Guildhall.Town)
	}

	guerra := "Não."
	if guild.Guild.Data.War {
		guerra = "Sim."
	}

	var membros []GuildMember
	for _, tipo := range guild.Guild.Members {
		for _, v := range tipo.Characters {
			var insert GuildMember
			insert.Name = v.Name
			insert.Nick = v.Nick
			insert.Level = v.Level
			insert.Vocation = v.Vocation
			insert.Status = v.Status
			membros = append(membros, insert)
		}
	}

	output := InternalGuild{}
	output.Name = guild.Guild.Data.Name
	output.Description = desc
	output.MemberCount = guild.Guild.Data.Totalmembers
	output.World = guild.Guild.Data.World
	output.GuildHall = guildHall
	output.War = guerra
	output.Members = membros

	return &output, nil
}

func CheckOnline(mundo string) ([]OnlineChar, *string, error) {
	if len(mundo) <= 0 {
		return nil, nil, errors.New("Você tem que especificar um mundo.")
	} else if len(mundo) < 2 {
		return nil, nil, ErrSmallName
	}

	world, err := GetWorld(mundo)
	if err != nil {
		return nil, nil, err
	}

	if len(world.World.WorldInformation.CreationDate) == 0 {
		return nil, nil, errors.New("Esse mundo não existe.")
	}

	var output []OnlineChar
	for _, v := range world.World.PlayersOnline {
		var insert OnlineChar
		insert.Name = v.Name
		insert.Level = v.Level
		insert.Vocation = v.Vocation
		output = append(output, insert)
	}

	return output, &world.World.WorldInformation.Name, nil
}

func GetTibiaNews(news ...int) (*InternalNews, error) {
	var inside int
	var url string

	switch len(news) {
	case 0:
		tibia, err := GetNews("news")
		if err != nil {
			return nil, err
		}
		inside = tibia.Newslist.Data[0].ID
		url = tibia.Newslist.Data[0].Tibiaurl
	case 1:
		inside = news[0]
		url = fmt.Sprintf("https://www.tibia.com/news/?subtopic=newsarchive&id=%d", inside)
	default:
		return nil, errors.New("getNews só aceita 1 argumento.")
	}

	tibiaInside, err := InsideNews(inside)
	if err != nil {
		return nil, err
	}

	out, err := formatNews(tibiaInside, url)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetTibiaNewsticker() (*InternalNews, error) {
	tibia, err := GetNews("ticker")
	if err != nil {
		return nil, err
	}

	url := tibia.Newslist.Data[0].Tibiaurl

	tibiaInside, err := InsideNews(tibia.Newslist.Data[0].ID)
	if err != nil {
		return nil, err
	}

	out, err := formatNews(tibiaInside, url)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func formatNews(tibiaInside *TibiaSpecificNews, url string) (*InternalNews, error) {
	if len(tibiaInside.News.Error) >= 1 {
		return nil, errors.New("Essa notícia não existe.")
	}

	t, err := dateparse.ParseLocal(tibiaInside.News.Date.Date)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`<(.*?)>`)
	desc := re.ReplaceAllString(tibiaInside.News.Content, "")
	shortdesc := ""

	if len(desc) > 1600 {
		split := strings.Split(desc, " ")
		for i := range split {
			if len(shortdesc) < 1600 {
				shortdesc += fmt.Sprintf("%s ", split[i])
			} else {
				shortdesc += "..."
				break
			}
		}
	} else {
		shortdesc = desc
	}

	output := InternalNews{}
	output.Title = tibiaInside.News.Title
	output.Description = desc
	output.ShortDescription = shortdesc
	output.URL = url
	output.Date = t.Format("02/01/2006")
	output.ID = tibiaInside.News.ID

	return &output, nil
}
