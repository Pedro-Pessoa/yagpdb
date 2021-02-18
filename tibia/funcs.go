package tibia

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/araddon/dateparse"
)

func GetTibiaChar(char string, update bool) (*InternalChar, error) {
	tibia, err := getChar(char)
	if err != nil {
		return nil, err
	}

	if tibia.Characters.Error != "" {
		return nil, errors.New(tibia.Characters.Error)
	}

	level := tibia.Characters.Data.Level
	if update {
		world, err := getWorld(tibia.Characters.Data.World)
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

	mortes := make([]InternalDeaths, 0, len(tibia.Characters.Deaths))
	for _, v := range tibia.Characters.Deaths {
		t2, err := dateparse.ParseLocal(v.Date.Date)
		if err != nil {
			return nil, err
		}

		mortes = append(mortes, InternalDeaths{
			Name:   tibia.Characters.Data.Name,
			Level:  v.Level,
			Reason: v.Reason,
			Date:   (t2.Add(time.Hour * -5)).Format("02/01/2006 15:04:05 BRT"),
		})
	}

	output := InternalChar{
		Name:              tibia.Characters.Data.Name,
		Level:             level,
		World:             tibia.Characters.Data.World,
		Vocation:          tibia.Characters.Data.Vocation,
		Residence:         tibia.Characters.Data.Residence,
		AccountStatus:     tibia.Characters.Data.AccountStatus,
		Status:            strings.Title(tibia.Characters.Data.Status),
		Loyalty:           lealdade,
		AchievementPoints: tibia.Characters.Data.AchievementPoints,
		Sex:               strings.Title(tibia.Characters.Data.Sex),
		Married:           casado,
		Guild:             guild,
		Rank:              cargo,
		Comment:           comentario,
		CreatedAt:         criado,
		House:             casa,
		Deaths:            mortes,
	}

	return &output, nil
}

func GetTibiaSpecificGuild(guildName string) (*InternalGuild, error) {
	err := validateName(guildName)
	if err != nil {
		return nil, err
	}

	guild, err := getSpecificGuild(strings.Title(guildName))
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
		guildHall = "**" + guild.Guild.Data.Guildhall.Name + "** que fica em " + guild.Guild.Data.Guildhall.Town
	}

	guerra := "Não."
	if guild.Guild.Data.War {
		guerra = "Sim."
	}

	membros := make([]GuildMember, 0, len(guild.Guild.Members))
	for _, tipo := range guild.Guild.Members {
		for _, v := range tipo.Characters {
			membros = append(membros, GuildMember{
				Name:     v.Name,
				Nick:     v.Nick,
				Level:    v.Level,
				Vocation: v.Vocation,
				Status:   v.Status,
			})
		}
	}

	output := InternalGuild{
		Name:        guild.Guild.Data.Name,
		Description: desc,
		MemberCount: guild.Guild.Data.Totalmembers,
		World:       guild.Guild.Data.World,
		GuildHall:   guildHall,
		War:         guerra,
		Members:     membros,
	}

	return &output, nil
}

func CheckOnline(mundo string) ([]OnlineChar, string, error) {
	if len(mundo) <= 0 {
		return nil, "", errors.New("Você tem que especificar um mundo.")
	} else if len(mundo) < 2 {
		return nil, "", ErrSmallName
	}

	world, err := getWorld(mundo)
	if err != nil {
		return nil, "", err
	}

	if len(world.World.WorldInformation.CreationDate) == 0 {
		return nil, "", errors.New("Esse mundo não existe.")
	}

	output := make([]OnlineChar, len(world.World.PlayersOnline))
	for i, v := range world.World.PlayersOnline {
		output[i] = OnlineChar(v)
	}

	return output, world.World.WorldInformation.Name, nil
}

func GetTibiaNews(news ...int) (*InternalNews, error) {
	var inside int
	var url string

	switch len(news) {
	case 0:
		tibia, err := getNews("news")
		if err != nil {
			return nil, err
		}
		inside = tibia.Newslist.Data[0].ID
		url = tibia.Newslist.Data[0].Tibiaurl
	case 1:
		inside = news[0]
		url = "https://www.tibia.com/news/?subtopic=newsarchive&id=" + strconv.Itoa(inside)
	default:
		return nil, errors.New("getNews só aceita 1 argumento.")
	}

	tibiaInside, err := insideNews(inside)
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
	tibia, err := getNews("ticker")
	if err != nil {
		return nil, err
	}

	url := tibia.Newslist.Data[0].Tibiaurl

	tibiaInside, err := insideNews(tibia.Newslist.Data[0].ID)
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
	var shortdesc strings.Builder

	if len(desc) > 1600 {
		split := strings.Split(desc, " ")
		for i := range split {
			if len(shortdesc.String()) < 1600 {
				shortdesc.WriteString(" " + split[i])
			} else {
				shortdesc.WriteString("...")
				break
			}
		}
	} else {
		shortdesc.WriteString(desc)
	}

	output := InternalNews{
		Title:            tibiaInside.News.Title,
		Description:      desc,
		ShortDescription: shortdesc.String(),
		URL:              url,
		Date:             t.Format("02/01/2006"),
		ID:               tibiaInside.News.ID,
	}

	return &output, nil
}
