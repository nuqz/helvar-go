package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type paramStringTestCase struct {
	p        Parameter
	expected string
}

func TestParameterString(t *testing.T) {
	testCases := map[string]paramStringTestCase{
		"string (Address)": {
			Parameter{ID: Address, Value: "x.y.z"},
			"@x.y.z",
		},
		"int": {
			Parameter{ID: Proportion, Value: int(-10)},
			"P:-10",
		},
		"uint": {
			Parameter{ID: Version, Value: uint(1)},
			"V:1",
		},
		"uint8": {
			Parameter{ID: Version, Value: uint8(1)},
			"V:1",
		},
		"uint16": {
			Parameter{ID: Version, Value: uint16(1)},
			"V:1",
		},
		"uint32": {
			Parameter{ID: Version, Value: uint32(1)},
			"V:1",
		},
		"uint64": {
			Parameter{ID: Mireds, Value: uint64(300)},
			"M:300",
		},
		"float32": {
			Parameter{ID: ColourX, Value: float32(0.77999)},
			"CX:0.78",
		},
		"float64": {
			Parameter{ID: ColourY, Value: float64(0.49222)},
			"CY:0.49",
		},
		// TODO: cover default case, when value type is not one of above
	}

	for tcDescription, tc := range testCases {
		t.Run(tcDescription, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.p.String())
		})
	}
}

type parseParamTestCase struct {
	pStr     string
	expected Parameter
}

func TestParseParameter(t *testing.T) {
	testCases := map[string]parseParamTestCase{
		"address": {
			"@1.2.3.4",
			Parameter{Address, "1.2.3.4"},
		},
		"negative proportion": {
			"P:-5",
			Parameter{Proportion, int64(-5)},
		},
		"integer": {
			"V:1",
			Parameter{Version, uint64(1)},
		},
		"float": {
			"CX:0.2499",
			Parameter{ColourX, float64(0.2499)},
		},
	}

	for tcDescription, tc := range testCases {
		t.Run(tcDescription, func(t *testing.T) {
			actual, err := ParseParameter(tc.pStr)
			require.NoError(t, err)
			assert.Equal(t, tc.expected.ID, actual.ID)
			assert.Equal(t, tc.expected.Value, actual.Value)
		})
	}
}
