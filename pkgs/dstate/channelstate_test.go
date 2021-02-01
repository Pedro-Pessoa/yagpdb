package dstate

import (
	"testing"
	"time"
)

func TestClearMessageBuffer(t *testing.T) {
	// test snowflakes: (2nd Sept. 2020)
	const flake3_21 = 750737145026707626 // 3:21 PM UTC
	const flake3_23 = 750737719482777743 // 3:23 PM UTC
	const flake3_28 = 750738872790679623 // 3:28 PM UTC
	const flake3_35 = 750740730032422992 // 3:35 PM UTC

	cases := []struct {
		name           string
		snowflakes     []int64
		maxMessageAge  time.Duration
		refTime        time.Time
		maxMessages    int
		expectedResult []int64
	}{
		{
			name: "max-messages-1",
			snowflakes: []int64{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			},
			maxMessages: 5,
			expectedResult: []int64{
				6, 7, 8, 9, 10,
			},
		},
		{
			name: "max-messages-2",
			snowflakes: []int64{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			},
			maxMessages: 10,
			expectedResult: []int64{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			},
		},
		{
			name: "max-message-age-1",
			snowflakes: []int64{
				flake3_21,
				flake3_23,
				flake3_28,
				flake3_35,
			},
			maxMessages:   10,
			maxMessageAge: time.Minute * 10,
			refTime:       createTestClearMessageBufferT(15, 37),
			expectedResult: []int64{
				flake3_28,
				flake3_35,
			},
		},
		{
			name: "max-message-age-2",
			snowflakes: []int64{
				flake3_21,
				flake3_23,
				flake3_28,
				flake3_35,
			},
			maxMessages:    10,
			maxMessageAge:  time.Minute * 10,
			refTime:        createTestClearMessageBufferT(16, 0),
			expectedResult: []int64{},
		},
		{
			name: "both-1",
			snowflakes: []int64{
				flake3_21,
				flake3_23,
				flake3_28,
				flake3_35,
			},
			maxMessages:   3,
			maxMessageAge: time.Minute * 10,
			refTime:       createTestClearMessageBufferT(15, 37),
			expectedResult: []int64{
				flake3_28,
				flake3_35,
			},
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			// create fake message states
			fakeInput := make([]*MessageState, len(v.snowflakes))
			for j, sf := range v.snowflakes {
				fakeInput[j] = &MessageState{
					ID: sf,
				}
			}

			output, _ := clearMessageBuffer(fakeInput, v.refTime, v.maxMessageAge, v.maxMessages)
			// create a intslice thats easier to print from the result
			intSlice := make([]int64, len(output))
			for j, ms := range output {
				intSlice[j] = ms.ID
			}

			if len(output) != len(v.expectedResult) {
				t.Errorf("Mismatched output lengths (%d vs expected %d) expected: %v, got %v", len(output), len(v.expectedResult), v.expectedResult, intSlice)
				return
			}

			// compare the full result content
			for j, id := range intSlice {
				if v.expectedResult[j] != id {
					t.Errorf("Unexpected output, expected: %v, got %v", v.expectedResult, intSlice)
					return
				}
			}
		})
	}
}

func createTestClearMessageBufferT(hour, min int) time.Time {
	t := time.Date(2020, time.September, 2, hour, min, 0, 0, time.UTC)
	return t
}
