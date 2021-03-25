package discordgo

import (
	"time"
)

// SnowflakeTimestamp returns the creation time of a Snowflake ID relative to the creation of Discord.
func SnowflakeTimestamp(ID int64) time.Time {
	timestamp := (ID >> 22) + 1420070400000
	return time.Unix(0, timestamp*1000000)
}

func CheckPerm(perms, permToCheck int64) bool {
	return permToCheck&perms == permToCheck
}

func AddPerms(perms int64, permsToAdd ...int64) int64 {
	for _, p := range permsToAdd {
		perms ^= p
	}

	return perms
}

func RemovePerms(perms int64, permsToRemove ...int64) int64 {
	for _, p := range permsToRemove {
		perms &^= p
	}

	return perms
}
