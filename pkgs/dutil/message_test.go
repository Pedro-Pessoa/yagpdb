package dutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StrSplitTestCase struct {
	Str      string
	ExpFirst string
	ExpRest  string
}

func TestSplitStr(t *testing.T) {
	cases := []StrSplitTestCase{
		{"123456789", "1234", "56789"},   // #0
		{"123\n456789", "123", "456789"}, // #1
		{"123 456789", "123", "456789"},  // #2
		{"1234", "1234", ""},             // #3
		{"123", "123", ""},               // #4
		{"12345", "1234", "5"},           // #5
		{"123 ", "123 ", ""},             // #6
		{"123 4", "123", "4"},            // #7
		{"123  ", "123", " "},            // #8
	}

	for k, c := range cases {
		first, last := StrSplit(c.Str, 4)
		assert.Equal(t, c.ExpFirst, first, "case #%d", k)
		assert.Equal(t, c.ExpRest, last, "case #%d", k)
	}
}

func TestSplitSendMessage(t *testing.T) {
	if !RequireTestingChannel(t) {
		return
	}

	msg := ""
	for i := 0; i < 400; i++ {
		msg += "TesMessage"
	}

	msgs, err := SplitSendMessage(dgo, envChannel, msg)
	if !assert.NoError(t, err, "Error sending messages") {
		return
	}
	assert.Len(t, msgs, 2, "Expected length is 2")
}
