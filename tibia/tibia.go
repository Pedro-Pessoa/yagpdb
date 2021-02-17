package tibia

import (
	"net/http"
	"strconv"

	"emperror.dev/errors"
	"github.com/mailru/easyjson"
)

func GetChar(name string) (*Tibia, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	resp, err := MakeRequest(name, "char")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tibia := Tibia{}
	err = easyjson.UnmarshalFromReader(resp.Body, &tibia)
	if err != nil {
		return nil, err
	}

	return &tibia, nil
}

func GetWorld(name string) (*TibiaWorld, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	valid := validWorld(name)
	if !valid {
		return nil, errors.New("O mundo " + name + " não existe")
	}

	resp, err := MakeRequest(name, "world")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	world := TibiaWorld{}
	err = easyjson.UnmarshalFromReader(resp.Body, &world)
	if err != nil {
		return nil, err
	}

	return &world, nil
}

func GetSpecificGuild(name string) (*SpecificGuild, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}

	resp, err := MakeRequest(name, "specificguild")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	specificGuild := SpecificGuild{}
	err = easyjson.UnmarshalFromReader(resp.Body, &specificGuild)
	if err != nil {
		return nil, err
	}

	return &specificGuild, nil
}

func GetNews(url string) (*TibiaNews, error) {
	resp, err := MakeRequest("news", url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tibia := TibiaNews{}
	err = easyjson.UnmarshalFromReader(resp.Body, &tibia)
	if err != nil {
		return nil, err
	}

	return &tibia, nil
}

func InsideNews(number int) (*TibiaSpecificNews, error) {
	resp, err := MakeRequest(number, "")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	tibiaInside := TibiaSpecificNews{}
	err = easyjson.UnmarshalFromReader(resp.Body, &tibiaInside)
	if err != nil {
		return nil, err
	}

	return &tibiaInside, nil
}

func MakeRequest(name interface{}, url string) (*http.Response, error) {
	var queryUrl string
	switch t := name.(type) {
	case string:
		switch url {
		case "specificguild":
			queryUrl = "https://api.tibiadata.com/v2/guild/" + t + ".json"
		case "news":
			queryUrl = "https://api.tibiadata.com/v2/latestnews.json"
		case "ticker":
			queryUrl = "https://api.tibiadata.com/v2/newstickers.json"
		case "world":
			queryUrl = "https://api.tibiadata.com/v2/world/" + t + ".json"
		default:
			queryUrl = "https://api.tibiadata.com/v2/characters/" + t + ".json"
		}
	case int:
		queryUrl = "https://api.tibiadata.com/v2/news/" + strconv.Itoa(t) + ".json"
	case int64:
		queryUrl = "https://api.tibiadata.com/v2/news/" + strconv.FormatInt(t, 10) + ".json"
	default:
		return nil, errors.New("Invalid name provided")
	}

	resp, err := http.DefaultClient.Get(queryUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "Erro no HTTP Get - MakeRequest Function")
	}

	return resp, nil
}
