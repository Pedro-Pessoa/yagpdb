package tibia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"emperror.dev/errors"
)

func GetChar(name string) (*Tibia, error) {
	if invalidName(name) {
		return nil, errors.New("O nome fornecido é inválido")
	}
	tibia := Tibia{}
	resp, err := MakeRequest(name, "char")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tibia)
	if err != nil {
		return nil, err
	}

	return &tibia, nil
}

func GetWorld(name string) (*TibiaWorld, error) {
	if invalidName(name) {
		return nil, errors.New("O nome fornecido é inválido")
	}
	world := TibiaWorld{}
	resp, err := MakeRequest(name, "world")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&world)
	if err != nil {
		return nil, err
	}

	return &world, nil
}

func GetSpecificGuild(name string) (*SpecificGuild, error) {
	if invalidName(name) {
		return nil, errors.New("O nome fornecido é inválido")
	}
	specificGuild := SpecificGuild{}
	resp, err := MakeRequest(name, "specificguild")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&specificGuild)
	if err != nil {
		return nil, err
	}

	return &specificGuild, nil
}

func GetNews(url string) (*TibiaNews, error) {
	tibia := TibiaNews{}
	resp, err := MakeRequest("news", url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tibia)
	if err != nil {
		return nil, err
	}

	return &tibia, nil
}

func InsideNews(number int) (*TibiaSpecificNews, error) {
	tibiaInside := TibiaSpecificNews{}
	resp, err := MakeRequest(number, "")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tibiaInside)
	if err != nil {
		return nil, err
	}

	return &tibiaInside, nil
}

func invalidName(name string) bool {
	matched, _ := regexp.MatchString(`[^\s'a-zA-Z\-\.]`, name)
	if matched {
		return true
	}
	return false
}

func MakeRequest(name interface{}, url string) (*http.Response, error) {
	var queryUrl string
	switch name.(type) {
	case string:
		switch url {
		case "specificguild":
			queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/guild/%s.json", name)
		case "news":
			queryUrl = "https://api.tibiadata.com/v2/latestnews.json"
		case "ticker":
			queryUrl = "https://api.tibiadata.com/v2/newstickers.json"
		case "world":
			queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/world/%s.json", name)
		default:
			queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/characters/%s.json", name)
		}
	case int, int64:
		queryUrl = fmt.Sprintf("https://api.tibiadata.com/v2/news/%d.json", name)
	default:
		return nil, nil
	}

	resp, err := http.DefaultClient.Get(queryUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "Erro no HTTP Get - MakeRequest Function")
	}

	return resp, nil
}
