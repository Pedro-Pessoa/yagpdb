package tibia

import (
	"bytes"
	"encoding/json"
)

type Tibia struct {
	Characters struct {
		Error string `json:"error"`
		Data  struct {
			Name              string   `json:"name"`
			FormerNames       []string `json:"former_names"`
			Title             string   `json:"title"`
			Sex               string   `json:"sex"`
			Vocation          string   `json:"vocation"`
			Level             int      `json:"level"`
			AchievementPoints int      `json:"achievement_points"`
			World             string   `json:"world"`
			FormerWorld       string   `json:"former_world"`
			Residence         string   `json:"residence"`
			MarriedTo         string   `json:"married_to"`
			House             struct {
				Name    string `json:"name"`
				Town    string `json:"town"`
				Paid    string `json:"paid"`
				World   string `json:"world"`
				Houseid int    `json:"houseid"`
			} `json:"house"`
			Guild struct {
				Name string `json:"name"`
				Rank string `json:"rank"`
			} `json:"guild"`
			LastLogin []struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"last_login"`
			Comment       string `json:"comment"`
			AccountStatus string `json:"account_status"`
			Status        string `json:"status"`
		} `json:"data"`
		Achievements []struct {
			Stars int    `json:"stars"`
			Name  string `json:"name"`
		} `json:"achievements"`
		Deaths []struct {
			Date struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
			Level    int    `json:"level"`
			Reason   string `json:"reason"`
			Involved []struct {
				Name string `json:"name"`
			} `json:"involved"`
		} `json:"deaths"`
		AccountInformation ActInfo `json:"account_information"`
		OtherCharacters    []struct {
			Name   string `json:"name"`
			World  string `json:"world"`
			Status string `json:"status"`
		} `json:"other_characters"`
	} `json:"characters"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type ActInfo struct {
	LoyaltyTitle string `json:"loyalty_title"`
	Created      struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"created"`
}

type TibiaWorld struct {
	World struct {
		WorldInformation struct {
			Name          string `json:"name"`
			PlayersOnline int    `json:"players_online"`
			OnlineRecord  struct {
				Players int `json:"players"`
				Date    struct {
					Date         string `json:"date"`
					TimezoneType int    `json:"timezone_type"`
					Timezone     string `json:"timezone"`
				} `json:"date"`
			} `json:"online_record"`
			CreationDate     string   `json:"creation_date"`
			Location         string   `json:"location"`
			PvpType          string   `json:"pvp_type"`
			WorldQuestTitles []string `json:"world_quest_titles"`
			BattleyeStatus   string   `json:"battleye_status"`
			GameWorldType    string   `json:"Game World Type:"`
		} `json:"world_information"`
		PlayersOnline []struct {
			Name     string `json:"name"`
			Level    int    `json:"level"`
			Vocation string `json:"vocation"`
		} `json:"players_online"`
	} `json:"world"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type TibiaNews struct {
	Newslist struct {
		Type string `json:"type"`
		Data []struct {
			ID       int    `json:"id"`
			Type     string `json:"type"`
			News     string `json:"news"`
			Apiurl   string `json:"apiurl"`
			Tibiaurl string `json:"tibiaurl"`
			Date     struct {
				Date         string `json:"date"`
				TimezoneType int    `json:"timezone_type"`
				Timezone     string `json:"timezone"`
			} `json:"date"`
		} `json:"data"`
	} `json:"newslist"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type TibiaSpecificNews struct {
	News struct {
		Error   string `json:"error"`
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Date    struct {
			Date         string `json:"date"`
			TimezoneType int    `json:"timezone_type"`
			Timezone     string `json:"timezone"`
		} `json:"date"`
	} `json:"news"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type SpecificGuild struct {
	Guild struct {
		Error string `json:"error"`
		Data  struct {
			Name          string     `json:"name"`
			Description   string     `json:"description"`
			Guildhall     GuildHouse `json:"guildhall"`
			Application   bool       `json:"application"`
			War           bool       `json:"war"`
			OnlineStatus  int        `json:"online_status"`
			OfflineStatus int        `json:"offline_status"`
			Disbanded     Finalizada `json:"disbanded"`
			Totalmembers  int        `json:"totalmembers"`
			Totalinvited  int        `json:"totalinvited"`
			World         string     `json:"world"`
			Founded       string     `json:"founded"`
			Active        bool       `json:"active"`
			Guildlogo     string     `json:"guildlogo"`
		} `json:"data"`
		Members []struct {
			RankTitle  string `json:"rank_title"`
			Characters []struct {
				Name     string `json:"name"`
				Nick     string `json:"nick"`
				Level    int    `json:"level"`
				Vocation string `json:"vocation"`
				Joined   string `json:"joined"`
				Status   string `json:"status"`
			} `json:"characters"`
		} `json:"members"`
		Invited []struct {
			Name    string `json:"name"`
			Invited string `json:"invited"`
		} `json:"invited"`
	} `json:"guild"`
	Information struct {
		APIVersion    int     `json:"api_version"`
		ExecutionTime float64 `json:"execution_time"`
		LastUpdated   string  `json:"last_updated"`
		Timestamp     string  `json:"timestamp"`
	} `json:"information"`
}

type Finalizada struct {
	Notification string `json:"notification"`
	Date         string `json:"date"`
}

type GuildHouse struct {
	Name    string `json:"name"`
	Town    string `json:"town"`
	Paid    string `json:"paid"`
	World   string `json:"world"`
	Houseid int    `json:"houseid"`
}

func (f *Finalizada) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("{")) {
		type finalizadaNoMethods Finalizada
		return json.Unmarshal(data, (*finalizadaNoMethods)(f))
	}
	return nil
}

func (gh *GuildHouse) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("{")) {
		type guildHouseNoMethods GuildHouse
		return json.Unmarshal(data, (*guildHouseNoMethods)(gh))
	}
	return nil
}

func (ai *ActInfo) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("{")) {
		type actInfoNoMethods ActInfo
		return json.Unmarshal(data, (*actInfoNoMethods)(ai))
	}
	return nil
}

type InternalChar struct {
	Name              string
	Level             int
	World             string
	Vocation          string
	Residence         string
	AccountStatus     string
	Status            string
	Loyalty           string
	AchievementPoints int
	Sex               string
	Married           string
	Guild             string
	Rank              string
	Comment           string
	CreatedAt         string
	House             string
	Deaths            []InternalDeaths
}

type InternalDeaths struct {
	Name   string
	Level  int
	Reason string
	Date   string
}

type InternalGuild struct {
	Name        string
	Description string
	MemberCount int
	World       string
	GuildHall   string
	War         string
	Members     []GuildMember
}

type GuildMember struct {
	Name     string
	Nick     string
	Level    int
	Vocation string
	Status   string
}

type OnlineChar struct {
	Name     string
	Level    int
	Vocation string
}

type InternalNews struct {
	Title            string
	Description      string
	ShortDescription string
	URL              string
	Date             string
	ID               int
}
