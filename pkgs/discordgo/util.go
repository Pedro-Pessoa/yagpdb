package discordgo

import (
	"time"
)

// SnowflakeTimestamp returns the creation time of a Snowflake ID relative to the creation of Discord.
func SnowflakeTimestamp(ID int64) time.Time {
	timestamp := (ID >> 22) + 1420070400000
	return time.Unix(0, timestamp*1000000)
}
