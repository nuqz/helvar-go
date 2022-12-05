package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type newCommandTestCase struct {
	version  uint8
	id       CommandID
	expected *Message
}

func TestNewCommand(t *testing.T) {
	testCases := map[string]newCommandTestCase{
		"control command - direct level device": {
			1, DirectLevelDevice,
			&Message{
				Type: TCommand,
				Parameters: []Parameter{
					{Version, uint8(1)},
					{Command, DirectLevelDevice},
				},
			},
		},
		// TODO: add more cases
	}

	for tcDescription, tc := range testCases {
		t.Run(tcDescription, func(t *testing.T) {
			msg := NewCommand(tc.version, tc.id)
			assert.Equal(t, tc.expected.Type, msg.Type)
			assert.Equal(t, tc.expected.Answer, msg.Answer)
			assert.Equal(t, tc.expected.IsPartial, msg.IsPartial)

			for i, param := range tc.expected.Parameters {
				assert.Equal(t, param.ID, msg.Parameters[i].ID)
				assert.Equal(t, param.Value, msg.Parameters[i].Value)
			}
		})
	}
}
