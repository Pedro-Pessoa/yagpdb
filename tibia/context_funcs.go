package tibia

import (
	"reflect"
	"sync"

	"emperror.dev/errors"

	"github.com/Pedro-Pessoa/tidbot/bot"
	"github.com/Pedro-Pessoa/tidbot/common/templates"
)

var ErrTooManyCalls = errors.New("Too many calls to this function")

func init() {
	templates.RegisterSetupFunc(func(c *templates.Context) {
		//Chars
		c.ContextFuncs["getChar"] = tmplGetTibiaChar(c)
		c.ContextFuncs["getDeaths"] = tmplGetCharDeaths(c)
		c.ContextFuncs["getDeath"] = tmplGetCharDeath(c)

		//Guild
		c.ContextFuncs["getGuild"] = tmplGetTibiaSpecificGuild(c)
		c.ContextFuncs["getGuildMembers"] = tmplGetTibiaSpecificGuildMembers(c)

		//Mundos
		c.ContextFuncs["checkOnline"] = tmplCheckOnline(c)

		//News
		c.ContextFuncs["getNews"] = tmplGetTibiaNews(c)
		c.ContextFuncs["getNewsticker"] = tmplGetTibiaNewsticker(c)

		//Tracks
		c.ContextFuncs["getTrackedHunteds"] = tmplGetTrackedHunteds(c)
		c.ContextFuncs["getTargetServerTracks"] = tmplGetTargetServerTracks(c)
		c.ContextFuncs["getTargetServerHunteds"] = tmplGetTargetServerHunteds(c)

		//Goroutine
		//Chars
		c.ContextFuncs["getMultipleChars"] = tmplGetMultipleChars(c)
		c.ContextFuncs["getMultipleCharsDeath"] = tmplGetMultipleCharsDeath(c)
	})
}

func structToMap(input interface{}) map[string]interface{} {
	v := reflect.ValueOf(input)
	typeOfS := v.Type()
	out := make(map[string]interface{})

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			out[typeOfS.Field(i).Name] = v.Field(i).Interface()
		}

		return out
	}

	return nil
}

func tmplGetTibiaChar(c *templates.Context) interface{} {
	return func(char string) (map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiachar", valueFree(c), 10) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaChar(char, true)
		if err != nil || output == nil {
			return nil, err
		}

		return structToMap(*output), nil
	}
}

func tmplGetCharDeaths(c *templates.Context) interface{} {
	return func(char string) ([]map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiachar", valueFree(c), 10) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaChar(char, false)
		if err != nil {
			return nil, err
		}

		outslice := make([]map[string]interface{}, len(output.Deaths))

		for k, v := range output.Deaths {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func tmplGetCharDeath(c *templates.Context) interface{} {
	return func(char string) (map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiachar", valueFree(c), 10) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaChar(char, false)
		if err != nil {
			return nil, err
		}

		if len(output.Deaths) == 0 {
			return nil, nil
		}

		return structToMap(output.Deaths[0]), nil
	}
}

func tmplGetTibiaSpecificGuild(c *templates.Context) interface{} {
	return func(guildName string) (map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiaguild", 1, 2) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaSpecificGuild(guildName)
		if err != nil || output == nil {
			return nil, err
		}

		return structToMap(*output), nil
	}
}

func tmplGetTibiaSpecificGuildMembers(c *templates.Context) interface{} {
	return func(guildName string) ([]map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiaguildmembers", 1, 2) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaSpecificGuild(guildName)
		if err != nil {
			return nil, err
		}

		outslice := make([]map[string]interface{}, len(output.Members))

		for k, v := range output.Members {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func tmplCheckOnline(c *templates.Context) interface{} {
	return func(mundo string) ([]map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibiamundo", 1, 2) {
			return nil, ErrTooManyCalls
		}

		output, _, err := CheckOnline(mundo)
		if err != nil {
			return nil, err
		}

		outslice := make([]map[string]interface{}, len(output))

		for k, v := range output {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func tmplGetTibiaNews(c *templates.Context) interface{} {
	return func(news ...int) (map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibianews", 1, 3) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaNews(news...)
		if err != nil || output == nil {
			return nil, err
		}

		return structToMap(*output), nil
	}
}

func tmplGetTibiaNewsticker(c *templates.Context) interface{} {
	return func() (map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("tibianews", 1, 3) {
			return nil, ErrTooManyCalls
		}

		output, err := GetTibiaNewsticker()
		if err != nil || output == nil {
			return nil, err
		}

		return structToMap(*output), nil
	}
}

func valueFree(c *templates.Context) int {
	memberCount := c.GS.Guild.MemberCount
	switch {
	case memberCount < 30:
		return 1
	case memberCount > 149:
		return 5
	default:
		return 3
	}
}

func tmplGetTrackedHunteds(c *templates.Context) interface{} {
	return func() ([]map[string]interface{}, error) {
		if c.IncreaseCheckCallCounterPremium("trackedhunteds", 1, 1) {
			return nil, ErrTooManyCalls
		}

		output, err := GetHuntedList(c.GS.ID)
		if err != nil {
			return nil, err
		}

		outslice := make([]map[string]interface{}, len(output))

		for k, v := range output {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func tmplGetTargetServerTracks(c *templates.Context) interface{} {
	return func(server int64) (interface{}, error) {
		if bot.IsGuildWhiteListed(c.GS.ID) {
			out, err := GetTracks(server)
			if err != nil {
				return nil, err
			}

			outslice := make([]map[string]interface{}, len(out))

			for k, v := range out {
				outslice[k] = structToMap(v)
			}

			return outslice, nil
		}

		return "", nil
	}
}

func tmplGetTargetServerHunteds(c *templates.Context) interface{} {
	return func(server int64) (interface{}, error) {
		if bot.IsGuildWhiteListed(c.GS.ID) {
			out, err := GetHuntedList(server)
			if err != nil {
				return nil, err
			}

			outslice := make([]map[string]interface{}, len(out))

			for k, v := range out {
				outslice[k] = structToMap(v)
			}

			return outslice, nil
		}

		return "", nil
	}
}

// Concurrent funcs

var wg sync.WaitGroup

func tmplGetMultipleChars(c *templates.Context) interface{} {
	return func(chars interface{}) ([]map[string]interface{}, error) {
		charsSlice, err := validateCharSlice(chars)
		if err != nil {
			return nil, err
		}

		output, err := GetTibiaMultiple(c, charsSlice, false)
		if err != nil {
			return nil, err
		}

		cast, ok := output.([]InternalChar)
		if !ok {
			logger.Error("Weird bug on tmplGetMultipleChars")
			return nil, nil
		}

		outslice := make([]map[string]interface{}, len(cast))

		for k, v := range cast {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func tmplGetMultipleCharsDeath(c *templates.Context) interface{} {
	return func(chars ...interface{}) ([]map[string]interface{}, error) {
		charsSlice, err := validateCharSlice(chars)
		if err != nil {
			return nil, err
		}

		output, err := GetTibiaMultiple(c, charsSlice, true)
		if err != nil {
			return nil, err
		}

		cast, ok := output.([]InternalDeaths)
		if !ok {
			logger.Error("Weird bug on tmplGetMultipleCharsDeath")
			return nil, nil
		}

		outslice := make([]map[string]interface{}, len(cast))

		for k, v := range cast {
			outslice[k] = structToMap(v)
		}

		return outslice, nil
	}
}

func validateCharSlice(input ...interface{}) ([]string, error) {
	var slice templates.Slice
	var err error

	switch l := len(input); l {
	case 0:
		return nil, errors.New("No arguments provided")
	case 1:
		switch t := input[0].(type) {
		case templates.Slice:
			slice = t
		case []interface{}:
			slice = t
		case string:
			return []string{t}, nil
		default:
			return nil, errors.Errorf("Invalid argument provided of type %T", t)
		}
	default:
		slice, err = templates.CreateSlice(input...)
		if err != nil {
			return nil, err
		}
	}

	if len(slice) > 10 {
		return nil, errors.Errorf("Você não pode solicitar mais do que 10 chars de uma vez")
	}

	out := make([]string, len(slice))
	for i, v := range slice {
		cast, ok := v.(string)
		if !ok {
			return nil, errors.Errorf("%v is not a string", v)
		}

		out[i] = cast
	}

	return out, nil
}

func GetTibiaMultiple(c *templates.Context, chars []string, deathsonly bool) (interface{}, error) {
	if c.IncreaseCheckCallCounterPremium("tibiagoroutine", 0, 1) {
		return nil, ErrTooManyCalls
	}

	fila := make(chan InternalChar, len(chars))
	for _, v := range chars {
		wg.Add(1)
		go charRoutine(fila, v)
	}

	wg.Wait()
	close(fila)

	if deathsonly {
		var output []InternalDeaths
		for e := range fila {
			if len(e.Deaths) > 0 {
				output = append(output, e.Deaths[0])
			}
		}

		return output, nil
	}

	var output []InternalChar
	for e := range fila {
		output = append(output, e)
	}

	return output, nil
}

func charRoutine(c chan InternalChar, char string) {
	defer wg.Done()

	ichar, err := GetTibiaChar(char, true)
	if err != nil || ichar == nil {
		return
	}

	c <- *ichar
}
