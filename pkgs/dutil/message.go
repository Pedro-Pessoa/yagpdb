package dutil

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Pedro-Pessoa/tidbot/pkgs/discordgo"
)

// SplitSendMessage see SplitSendMessageCtx
func SplitSendMessage(s *discordgo.Session, channelID int64, message string) ([]*discordgo.Message, error) {
	return SplitSendMessageCtx(s, context.Background(), channelID, message)
}

// SplitSendMessageCtx is a helper for sending potentially long messages
// If the message is longer than 2k characters it will split at
// Last newline before 2k or last whitespace before 2k or if that fails
// (no whitespace) just split at 2k
func SplitSendMessageCtx(s *discordgo.Session, ctx context.Context, channelID int64, message string) ([]*discordgo.Message, error) {
	return SplitSendMessagePSCtx(s, ctx, channelID, message, "", "", false, false)
}

// SplitSendMessagePS see SplitSendMessagePSCtx
func SplitSendMessagePS(s *discordgo.Session, channelID int64, message string, prefix, suffix string, prefixStart, suffixEnd bool) ([]*discordgo.Message, error) {
	return SplitSendMessagePSCtx(s, context.Background(), channelID, message, prefix, suffix, prefixStart, suffixEnd)
}

// SplitSendMessagePSCtx is a helper for sending potentially long messages
// If the message is longer than 2k characters it will split at
// Last newline before 2k or last whitespace before 2k or if that fails
// (no whitespace) just split at 2k
// Prefix is added to the start of each message sent (usefull for codeblocks),
// Prefix is not not added to the first one if prefixStart is false
// Suffix is added to the end of each message, and not the last message if suffixend is false
// Cancel the context to stop this process
func SplitSendMessagePSCtx(s *discordgo.Session, ctx context.Context, channelID int64, message string, prefix, suffix string, prefixStart, suffixEnd bool) ([]*discordgo.Message, error) {
	rest := message
	first := true

	ret := make([]*discordgo.Message, 0)

	for {

		if ctx.Err() != nil {
			return ret, ctx.Err()
		}

		maxLen := 2000

		// Take away prefix and suffix length if used
		if prefixStart || !first {
			maxLen -= utf8.RuneCountInString(prefix)
		}
		maxLen -= utf8.RuneCountInString(suffix)

		msg, newRest := StrSplit(rest, maxLen)

		// Add the actual prefix and suffix
		if prefixStart || !first {
			msg = prefix + msg
		}
		if suffixEnd || len(newRest) > 0 {
			msg += suffix
		}

		discordMessage, err := s.ChannelMessageSend(channelID, msg)
		if err != nil {
			return nil, err
		}

		ret = append(ret, discordMessage)

		rest = newRest
		if rest == "" {
			break
		}

		first = false
	}

	return ret, nil
}

// Will split "s" before runecount at last possible newline, whitespace or just at "runecount" if there is no whitespace
// If the runecount in "s" is less than "runeCount" then "last" will be zero
func StrSplit(s string, runeCount int) (split, rest string) {
	// Possibly split up s
	if utf8.RuneCountInString(s) > runeCount {
		_, beforeIndex := RuneByIndex(s, runeCount)
		firstPart := s[:beforeIndex]

		// Split at newline if possible
		foundWhiteSpace := false
		lastIndex := strings.LastIndex(firstPart, "\n")
		if lastIndex == -1 {
			// No newline, check for any possible whitespace then
			lastIndex = strings.LastIndexFunc(firstPart, func(r rune) bool {
				return unicode.In(r, unicode.White_Space)
			})
			if lastIndex == -1 {
				lastIndex = beforeIndex
			} else {
				foundWhiteSpace = true
			}
		} else {
			foundWhiteSpace = true
		}

		// Remove the whitespace we split at if any
		if foundWhiteSpace {
			_, rLen := utf8.DecodeRuneInString(s[lastIndex:])
			rest = s[lastIndex+rLen:]
		} else {
			rest = s[lastIndex:]
		}

		split = s[:lastIndex]
	} else {
		split = s
	}

	return
}

// Returns the string index from the rune position
// Panics if utf8.RuneCountInString(s) <= runeIndex or runePos < 0
func RuneByIndex(s string, runePos int) (rune, int) {
	sLen := utf8.RuneCountInString(s)
	if sLen <= runePos || runePos < 0 {
		panic("runePos is out of bounds")
	}

	i := 0
	last := rune(0)
	for k, r := range s {
		if i == runePos {
			return r, k
		}
		i++
		last = r
	}
	return last, i
}
