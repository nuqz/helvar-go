package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ParameterID string

const (
	Version            ParameterID = "V"
	Command            ParameterID = "C"
	Address            ParameterID = "@"
	Group              ParameterID = "G"
	Block              ParameterID = "B"
	Scene              ParameterID = "S"
	FadeTime           ParameterID = "F"
	Level              ParameterID = "L"
	Proportion         ParameterID = "P"
	DisplayScreen      ParameterID = "D"
	SequenceNumber     ParameterID = "Q"
	Time               ParameterID = "T"
	Ack                ParameterID = "A"
	Latitude           ParameterID = "L"
	Longitude          ParameterID = "E"
	TimeZoneDifference ParameterID = "Z"
	DaylightSavingTime ParameterID = "Y"
	ConstantLightScene ParameterID = "K"
	ForceStoreScene    ParameterID = "O"
	Mireds             ParameterID = "M"
	ColourX            ParameterID = "CX"
	ColourY            ParameterID = "CY"
)

// Parameter represents a pair of message parameter ID and message
// parameter value.
type Parameter struct {
	ID    ParameterID
	Value any
}

// ParseParameter returns message parameter when input is valid string
// representation, otherwise it returns an error.
func ParseParameter(input string) (Parameter, error) {
	if string(input[0]) == string(Address) {
		return Parameter{
			ID:    Address,
			Value: input[1:],
		}, nil
	}

	idValue := strings.Split(input, ParameterIDDelimeter.String())
	if len(idValue) != 2 {
		return Parameter{}, errors.Errorf(
			`"%s" is not valid message parameter string`, input)
	}

	param := Parameter{ID: ParameterID(idValue[0])}
	if vUint, err := strconv.ParseUint(idValue[1], 10, 0); err == nil {
		param.Value = vUint
	} else {
		if vInt, err := strconv.ParseInt(idValue[1], 10, 64); err == nil {
			param.Value = vInt
		} else {
			if vFloat, err := strconv.ParseFloat(idValue[1], 64); err == nil {
				param.Value = vFloat
			} else {
				param.Value = idValue[1]
			}
		}
	}

	if param.ID == Command {
		param.Value = CommandID(param.Value.(uint64))
	}

	return param, nil
}

// String returns message parameter serialized in string form.
func (p Parameter) String() string {
	// Exception for Address parameter
	if p.ID == Address {
		return fmt.Sprintf("%s%s", p.ID, p.Value)
	}

	format := "%s%s%v"
	switch p.Value.(type) {
	case string:
		format = "%s%s%s"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		format = "%s%s%d"
	case float32, float64:
		format = "%s%s%.2f"
	}

	return fmt.Sprintf(format, p.ID, ParameterIDDelimeter, p.Value)
}
