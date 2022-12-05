package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parseTestCase struct {
	msgTxt   string
	expected *Message
}

func TestParse(t *testing.T) {
	testCases := map[string]parseTestCase{
		"control command - recall scene": {
			">V:1,C:11,G:1,S:1#",
			&Message{
				Type: TCommand,
				Parameters: []Parameter{
					{Version, uint64(1)},
					{Command, RecallSceneGroup},
					{Group, uint64(1)},
					{Scene, uint64(1)},
				},
				Answer:    "",
				IsPartial: false,
			},
		},
		"query reply - query groups": {
			"?V:1,C:165=1,2,3,4,5#",
			&Message{
				Type: TReply,
				Parameters: []Parameter{
					{Version, uint64(1)},
					{Command, QueryGroups},
				},
				Answer:    "1,2,3,4,5",
				IsPartial: false,
			},
		},
		"partial query reply - query groups": {
			"?V:1,C:165=5,4,3,2,1$",
			&Message{
				Type: TReply,
				Parameters: []Parameter{
					{Version, uint64(1)},
					{Command, QueryGroups},
				},
				Answer:    "5,4,3,2,1",
				IsPartial: true,
			},
		},
		"error query reply - query group description": {
			"!V:1,C:105,G:9999=1#",
			&Message{
				Type: TError,
				Parameters: []Parameter{
					{Version, uint64(1)},
					{Command, QueryGroupDescription},
					{Group, uint64(9999)},
				},
				Answer:    "1",
				IsPartial: false,
			},
		},
	}

	for tcDescription, tc := range testCases {
		t.Run(tcDescription, func(t *testing.T) {
			actual, err := Parse(tc.msgTxt)
			require.NoError(t, err)

			assert.Equal(t, tc.expected.Type, actual.Type)
			assert.Equal(t, tc.expected.Answer, actual.Answer)
			assert.Equal(t, tc.expected.IsPartial, actual.IsPartial)

			require.Equal(t,
				len(tc.expected.Parameters), len(actual.Parameters))
			for i, param := range tc.expected.Parameters {
				assert.Equal(t, param.ID, actual.Parameters[i].ID)
				assert.Equal(t, param.Value, actual.Parameters[i].Value, param.ID)
			}
		})
	}
}

type msgStringTestCase struct {
	msg      *Message
	expected string
}

func TestMessageToString(t *testing.T) {
	testCases := map[string]msgStringTestCase{
		"control command - recall scene": {
			&Message{
				Type: TCommand,
				Parameters: []Parameter{
					{Version, 1},
					{Command, RecallSceneGroup},
					{Group, 2},
					{Block, 3},
					{Scene, 4},
				},
				Answer:    "",
				IsPartial: false,
			},
			">V:1,C:11,G:2,B:3,S:4#",
		},
		// TODO:
		// - Direct level (group)
		// - Direct level (device)
		// - Direct level (colour)
		// - Direct level (temperature/mireds)
		// - Query commands
		// - Query replies
		// - Internal
		// - Error/diagnostics
	}

	for tcDescription, tc := range testCases {
		t.Run(tcDescription, func(t *testing.T) {
			actual := tc.msg.String()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
