package tibia

// API Structs
type Tibia struct {
	Characters  TibiaCharacters `json:"characters"`
	Information ApiInformation  `json:"information"`
}

type TibiaCharacters struct {
	Error              string             `json:"error"`
	Data               CharData           `json:"data"`
	Achievements       []CharAchievements `json:"achievements"`
	Deaths             []CharDeaths       `json:"deaths"`
	AccountInformation ActInfo            `json:"account_information"`
	OtherCharacters    []OtherChars       `json:"other_characters"`
}

type CharData struct {
	Name              string          `json:"name"`
	FormerNames       []string        `json:"former_names"`
	Title             string          `json:"title"`
	Sex               string          `json:"sex"`
	Vocation          string          `json:"vocation"`
	Level             int             `json:"level"`
	AchievementPoints int             `json:"achievement_points"`
	World             string          `json:"world"`
	FormerWorld       string          `json:"former_world"`
	Residence         string          `json:"residence"`
	MarriedTo         string          `json:"married_to"`
	House             CharHouse       `json:"house"`
	Guild             CharGuild       `json:"guild"`
	LastLogin         []CharLastLogin `json:"last_login"`
	Comment           string          `json:"comment"`
	AccountStatus     string          `json:"account_status"`
	Status            string          `json:"status"`
}

type CharHouse struct {
	Name    string `json:"name"`
	Town    string `json:"town"`
	Paid    string `json:"paid"`
	World   string `json:"world"`
	Houseid int    `json:"houseid"`
}

type CharGuild struct {
	Name string `json:"name"`
	Rank string `json:"rank"`
}

type CharLastLogin struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type CharAchievements struct {
	Stars int    `json:"stars"`
	Name  string `json:"name"`
}

type CharDeaths struct {
	Date     Date            `json:"date"`
	Level    int             `json:"level"`
	Reason   string          `json:"reason"`
	Involved []DeathInvolved `json:"involved"`
}

type DeathInvolved struct {
	Name string `json:"name"`
}

type OtherChars struct {
	Name   string `json:"name"`
	World  string `json:"world"`
	Status string `json:"status"`
}

type ActInfo struct {
	LoyaltyTitle string     `json:"loyalty_title"`
	Created      ActCreated `json:"created"`
}

type ActCreated struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type TibiaWorld struct {
	World       World          `json:"world"`
	Information ApiInformation `json:"information"`
}

type World struct {
	WorldInformation WorldInformation `json:"world_information"`
	PlayersOnline    []PlayersOnline  `json:"players_online"`
}

type WorldInformation struct {
	Name             string       `json:"name"`
	PlayersOnline    int          `json:"players_online"`
	OnlineRecord     OnlineRecord `json:"online_record"`
	CreationDate     string       `json:"creation_date"`
	Location         string       `json:"location"`
	PvpType          string       `json:"pvp_type"`
	WorldQuestTitles []string     `json:"world_quest_titles"`
	BattleyeStatus   string       `json:"battleye_status"`
	GameWorldType    string       `json:"Game World Type:"`
}

type OnlineRecord struct {
	Players int  `json:"players"`
	Date    Date `json:"date"`
}

type PlayersOnline struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
}

type TibiaNews struct {
	Newslist    Newslist       `json:"newslist"`
	Information ApiInformation `json:"information"`
}

type Newslist struct {
	Type string     `json:"type"`
	Data []NewsData `json:"data"`
}

type NewsData struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	News     string `json:"news"`
	Apiurl   string `json:"apiurl"`
	Tibiaurl string `json:"tibiaurl"`
	Date     Date   `json:"date"`
}

type TibiaSpecificNews struct {
	News        News           `json:"news"`
	Information ApiInformation `json:"information"`
}

type News struct {
	Error   string `json:"error"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    Date   `json:"date"`
}

type SpecificGuild struct {
	Guild       Guild          `json:"guild"`
	Information ApiInformation `json:"information"`
}

type Guild struct {
	Error   string         `json:"error"`
	Data    GuildData      `json:"data"`
	Members []GuildMembers `json:"members"`
	Invited []GuildInvited `json:"invited"`
}

type GuildMembers struct {
	RankTitle  string                  `json:"rank_title"`
	Characters []GuildMemberCharacters `json:"characters"`
}

type GuildMemberCharacters struct {
	Name     string `json:"name"`
	Nick     string `json:"nick"`
	Level    int    `json:"level"`
	Vocation string `json:"vocation"`
	Joined   string `json:"joined"`
	Status   string `json:"status"`
}

type GuildInvited struct {
	Name    string `json:"name"`
	Invited string `json:"invited"`
}

type GuildData struct {
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
}

type GuildHouse struct {
	Name    string `json:"name"`
	Town    string `json:"town"`
	Paid    string `json:"paid"`
	World   string `json:"world"`
	Houseid int    `json:"houseid"`
}

type Finalizada struct {
	Notification string `json:"notification"`
	Date         string `json:"date"`
}

type Date struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type ApiInformation struct {
	APIVersion    int     `json:"api_version"`
	ExecutionTime float64 `json:"execution_time"`
	LastUpdated   string  `json:"last_updated"`
	Timestamp     string  `json:"timestamp"`
}

// Internal Structs
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
