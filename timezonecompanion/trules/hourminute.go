package trules

import (
	"regexp"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/jonas747/when/rules"
)

func HourMinute(s rules.Strategy) rules.Rule {
	return &rules.F{
		RegExp: regexp.MustCompile(`(?i)(?:\W|^)((?:[0-1]{0,1}[0-9])|(?:2[0-3]))(?:\:|：)((?:[0-5][0-9]))(?:\s*(A\.M\.|P\.M\.|AM|PM))?(?:\W|$)`),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			if (c.Hour != nil || c.Minute != nil) && s != rules.Override {
				return false, nil
			}

			hour, err := strconv.Atoi(m.Captures[0])
			if err != nil {
				return false, errors.WrapIf(err, "hour minute rule")
			}

			minutes, err := strconv.Atoi(m.Captures[1])
			if err != nil {
				return false, errors.WrapIf(err, "hour minute rule")
			}

			if minutes > 59 {
				return false, nil
			}
			c.Minute = &minutes

			if m.Captures[2] != "" {
				if hour > 12 {
					return false, nil
				}
				switch m.Captures[2][0] {
				case 65, 97: // am
					if hour == 12 {
						hour -= 12
					}
					c.Hour = &hour
				case 80, 112: // pm
					if hour < 12 {
						hour += 12
					}
					c.Hour = &hour
				}
			} else {
				if hour > 23 {
					return false, nil
				}
				c.Hour = &hour
			}

			return true, nil
		},
	}
}
