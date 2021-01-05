package tibia

import (
	"reflect"
	"sync"

	"emperror.dev/errors"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/common/templates"
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
		if err != nil {
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
		if err != nil {
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
		if err != nil {
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
		if err != nil {
			return nil, err
		}

		return structToMap(*output), nil
	}
}

func valueFree(c *templates.Context) int {
	memberCount := c.GS.Guild.MemberCount
	valueFree := 3
	if memberCount < 30 {
		valueFree = 1
	} else if memberCount > 149 {
		valueFree = 5
	}
	return valueFree
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

//////////// Go Routine Stuff ////////////

var wg sync.WaitGroup

func tmplGetMultipleChars(c *templates.Context) interface{} {
	return func(chars interface{}) ([]map[string]interface{}, error) {
		output, err := GetTibiaMultiple(c, chars, false)
		if err != nil {
			return nil, err
		}

		switch t := output.(type) {
		case []InternalChar:
			if len(t) > 0 {
				outslice := make([]map[string]interface{}, len(t))
				for k, v := range t {
					outslice[k] = structToMap(v)
				}

				return outslice, nil
			}
		}

		return nil, nil
	}
}

func tmplGetMultipleCharsDeath(c *templates.Context) interface{} {
	return func(chars interface{}) ([]map[string]interface{}, error) {
		output, err := GetTibiaMultiple(c, chars, true)
		if err != nil {
			return nil, err
		}

		switch t := output.(type) {
		case []InternalDeaths:
			if len(t) > 0 {
				outslice := make([]map[string]interface{}, len(t))
				for k, v := range t {
					outslice[k] = structToMap(v)
				}

				return outslice, nil
			}
		}

		return nil, nil
	}
}

func GetTibiaMultiple(c *templates.Context, chars interface{}, deathsonly bool) (interface{}, error) {
	if c.IncreaseCheckCallCounterPremium("tibiagoroutine", 0, 1) {
		return nil, ErrTooManyCalls
	}

	v := reflect.ValueOf(chars)
	var slice []interface{}
	switch v.Kind() {
	case reflect.Slice:
		if v.Len() > 10 {
			return nil, errors.New("você não pode solicitar mais do que 10 personagens de uma vez.")
		}
		switch t := chars.(type) {
		case []interface{}:
			slice = t
		case templates.Slice:
			slice = t
		}
	default:
		return nil, errors.New("Essa função só aceita slices como argumento.")
	}

	fila := make(chan InternalChar, len(slice))
	for _, v := range slice {
		switch t := v.(type) {
		case string:
			wg.Add(1)
			go charRoutine(fila, t)
		default:
			continue
		}
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
	defer routineCleanUp()
	ichar, err := GetTibiaChar(char, true)
	if err != nil {
		return
	}
	c <- *ichar
	return
}

func routineCleanUp() {
	if r := recover(); r != nil {
		logger.Infof("Recovered at: %v", r)
	}
	defer wg.Done()
}
