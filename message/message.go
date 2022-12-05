package message

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Char byte

func (c Char) String() string { return string(c) }
func (c Char) Byte() byte     { return byte(c) }

const (
	TCommand               Char = '>'
	TInternalCommand       Char = '<'
	TReply                 Char = '?'
	TError                 Char = '!'
	Delimiter              Char = ','
	AddressDelimiter       Char = '.'
	ParameterIDDelimeter   Char = ':'
	Terminator             Char = '#'
	PartialReplyTerminator Char = '$'
	Answer                 Char = '='

	// How it looks like:
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	MaxMessageBytes = 1500
)

var (
	allowedStartChars = []Char{
		TCommand,
		TInternalCommand,
		TReply,
		TError,
	}
	allowedEndChars = []Char{
		Terminator,
		PartialReplyTerminator,
	}
)

// Message is a base for both commands and replies. We'll use only ASCII
// format, raw binary format will not be covered by this package.
//
// The following is from HelvarNET protocol docs.
// Any message sent to, or received from, a router can be in either ASCII or
// raw binary form (see Command Format for more information).
// The format of the data contained within messages is defined by the
// protocol. A query reply message from the router will be in the same format
// as the query command message sent i.e. if a query message is sent in ASCII
// form then the reply will also be in ASCII.
//
// Messages must not exceed the maximum length of 1500 bytes.
// TODO: validate ^^^ before send.
type Message struct {
	Type       Char
	Parameters []Parameter
	Answer     string
	IsPartial  bool
}

func (msg *Message) SetType(t Char) *Message {
	msg.Type = t
	return msg
}

// Parse returns a Message parsed from a given input string (HelvarNET
// ASCII format).
//
// General format:
// <MessageType><MessageParameter[, ...]>=<MessageResult><#|$>
func Parse(input string) (*Message, error) {
	start, end := Char(input[0]), Char(input[len(input)-1])
	if !slices.Contains(allowedStartChars, start) ||
		!slices.Contains(allowedEndChars, end) {
		return nil, errors.Errorf("failed to parse message: %s", input)
	}

	body := strings.Trim(input, string(start)+end.String())
	bodyParts := strings.Split(body, Answer.String())
	result := ""
	if len(bodyParts) == 2 {
		result = bodyParts[1]
	}

	rawParams := strings.Split(bodyParts[0], Delimiter.String())
	params := make([]Parameter, len(rawParams))

	var err error
	for i, rp := range rawParams {
		params[i], err = ParseParameter(rp)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:       start,
		Parameters: params,
		Answer:     result,
		IsPartial:  end == PartialReplyTerminator,
	}, nil
}

func ParsePartial(in string) (*Message, error) {
	strs := strings.SplitAfter(in, PartialReplyTerminator.String())
	ln := len(strs)

	out := &Message{}
	answers := make([]string, ln)
	for i, str := range strs {
		msg, err := Parse(str)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			out.Type = msg.Type
			out.Parameters = slices.Clone(msg.Parameters)
			out.IsPartial = ln > 1
		}

		answers[i] = msg.Answer
	}

	out.Answer = strings.Join(answers, Delimiter.String())
	return out, nil
}

func (msg *Message) ParamsToString() string {
	strParams := make([]string, len(msg.Parameters))
	for i, param := range msg.Parameters {
		strParams[i] = param.String()
	}
	return strings.Join(strParams, Delimiter.String())
}

func (msg *Message) ID() string {
	return msg.ParamsToString()
}

// String returns a Message represented as string (HelvarNET ASCII format).
func (msg *Message) String() string {
	msgBody := msg.ParamsToString()
	if msg.Answer != "" {
		msgBody += Answer.String() + msg.Answer
	}

	return string(msg.Type) + msgBody + string(Terminator)
}

func (msg *Message) Bytes() []byte {
	return []byte(msg.String())
}

// AddParameters adds parameters to a given message and returns message itself
// to support methods chaining.
func (msg *Message) AddParameters(params ...Parameter) *Message {
	msg.Parameters = append(msg.Parameters, params...)
	return msg
}

// GetParameter returns parameter value requested by its ID.
func (msg *Message) GetParameter(id ParameterID) any {
	for _, param := range msg.Parameters {
		if param.ID == id {
			return param.Value
		}
	}

	return nil
}

// GetCommandID returns command ID of a given message if command parameter is
// present, otherwise it returns 0 (NoCommand).
func (msg *Message) GetCommandID() CommandID {
	if v := msg.GetParameter(Command); v != nil {
		return v.(CommandID)
	}

	return NoCommand
}

func (msg *Message) GetGroupID() uint16 {
	if i := msg.GetParameter(Group); i != nil {
		switch v := i.(type) {
		case uint16:
			return v
		case uint64:
			return uint16(v)
		}
	}

	return 0
}

func (msg *Message) GetAddress() string {
	if i := msg.GetParameter(Address); i != nil {
		return i.(string)
	}

	return ""
}

func (msg *Message) AnswerStrings() []string {
	return strings.Split(msg.Answer, Delimiter.String())
}

func (msg *Message) AnswerIDs() ([]int, error) {
	strs := msg.AnswerStrings()

	out := []int{}
	for _, idStr := range strs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}

	return out, nil
}
