package util

import (
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

func ParseColor(raw string) (int, bool) {
	raw = strings.TrimPrefix(raw, "#")

	// try to parse as hex color code first
	parsed, err := strconv.ParseInt(raw, 16, 32)
	if err == nil {
		return int(parsed), true
	}

	// look up the color code table
	for _, v := range colornames.Names {
		if strings.EqualFold(v, raw) {
			cStruct := colornames.Map[v]

			color := (int(cStruct.R) << 16) | (int(cStruct.G) << 8) | int(cStruct.B)
			return color, true
		}
	}

	return 0, false
}
