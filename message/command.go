package message

import (
	"image/color"
	"math"

	"github.com/nuqz/col2xy"
	"golang.org/x/exp/slices"
)

type CommandID uint8

const (
	// Package internals

	NoCommand CommandID = 0

	// Control

	RecallSceneGroup  CommandID = 11
	RecallSceneDevice CommandID = 12
	DirectLevelGroup  CommandID = 13
	DirectLevelDevice CommandID = 14

	// Query

	QueryClusters                CommandID = 101
	QueryRouters                 CommandID = 102
	QueryGroupDescription        CommandID = 105
	QueryDeviceDescription       CommandID = 106
	QueryDeviceTypesAndAddresses CommandID = 100
	QueryDeviceState             CommandID = 110
	QueryWorkgroupName           CommandID = 107
	QueryDeviceLoadLevel         CommandID = 152
	QuerySceneInfo               CommandID = 167
	QueryTime                    CommandID = 185
	QueryLastSceneInGroup        CommandID = 109
	QueryLastSceneInBlock        CommandID = 103
	QueryGroup                   CommandID = 164
	QueryGroups                  CommandID = 165
	QuerySceneNames              CommandID = 166
	QueryRouterVersion           CommandID = 190
	QueryHelvarnetVersion        CommandID = 191
)

var (
	CommandsWithoutResponse = []CommandID{
		RecallSceneGroup,
		RecallSceneDevice,
		DirectLevelGroup,
		DirectLevelDevice,
	}
)

func NeedResponse(msg *Message) bool {
	return !slices.Contains(CommandsWithoutResponse, msg.GetCommandID())
}

func NewCommand(version uint8, id CommandID) *Message {
	out := &Message{}
	return out.SetType(TCommand).
		AddParameters(
			Parameter{Version, version},
			Parameter{Command, id},
		)
}

func NewCommandV1(id CommandID) *Message { return NewCommand(1, id) }

func NewQueryRouters(address string) *Message {
	return NewCommandV1(QueryRouters).
		AddParameters(Parameter{Address, address})
}

func NewQueryTime() *Message { return NewCommandV1(QueryTime) }

func NewQueryClusters() *Message { return NewCommandV1(QueryClusters) }

func NewQueryGroups() *Message { return NewCommandV1(QueryGroups) }

func NewQueryGroupDescription(gid uint16) *Message {
	return NewCommandV1(QueryGroupDescription).
		AddParameters(Parameter{Group, gid})
}

func NewQueryGroup(gid uint16) *Message {
	return NewCommandV1(QueryGroup).
		AddParameters(Parameter{Group, gid})
}

func NewQueryDeviceDescription(address string) *Message {
	return NewCommandV1(QueryDeviceDescription).
		AddParameters(Parameter{Address, address})
}

func NewRecallScene(
	cmdID CommandID,
	block, scene uint8,
	params ...Parameter,
) *Message {
	return NewCommandV1(cmdID).AddParameters(
		Parameter{Block, block},
		Parameter{Scene, scene}).
		AddParameters(params...)
}

func NewRecallSceneGroup(
	gid uint16,
	block, scene uint8,
	params ...Parameter,
) *Message {
	return NewRecallScene(RecallSceneGroup, block, scene, params...).
		AddParameters(Parameter{Group, gid})
}

func NewRecallSceneDevice(
	addr string,
	block, scene uint8,
	params ...Parameter,
) *Message {
	return NewRecallScene(RecallSceneDevice, block, scene, params...).
		AddParameters(Parameter{Address, addr})
}

func NewDirectLevel(
	cmdID CommandID,
	level uint8,
	params ...Parameter,
) *Message {
	return NewCommandV1(cmdID).
		AddParameters(Parameter{Level, level}).
		AddParameters(params...)
}

func NewDirectLevelGroup(
	gid uint16,
	level uint8,
	params ...Parameter,
) *Message {
	return NewDirectLevel(DirectLevelGroup, level, params...).
		AddParameters(Parameter{Group, gid})
}

func NewDirectLevelDevice(
	addr string,
	level uint8,
	params ...Parameter,
) *Message {
	return NewDirectLevel(DirectLevelDevice, level, params...).
		AddParameters(Parameter{Address, addr})
}

func Kelvins2Mireds(k uint16) int {
	// The formula was taken from one of Helvar manuals
	return int(math.Round(1000000 / float64(k)))
}

func Kelvins2MiredsParam(k uint16) Parameter {
	return Parameter{Mireds, Kelvins2Mireds(k)}
}

func NewColorTemperatureGroup(
	gid, tempK uint16,
	level uint8,
	params ...Parameter,
) *Message {
	return NewDirectLevelGroup(gid, level, params...).
		AddParameters(Kelvins2MiredsParam(tempK))
}

func NewColorTemperatureDevice(
	addr string,
	tempK uint16,
	level uint8,
	params ...Parameter,
) *Message {
	return NewDirectLevelDevice(addr, level, params...).
		AddParameters(Kelvins2MiredsParam(tempK))
}

// TODO: Other shortcuts for:
// - proportions
// - other control commands...

func NewCommandV2(id CommandID) *Message { return NewCommand(2, id) }

func newColor(
	cmdID CommandID,
	level uint8,
	x, y float64,
	params ...Parameter,
) *Message {
	return NewCommandV2(cmdID).
		AddParameters(
			Parameter{Level, level},
			Parameter{ColourX, x},
			Parameter{ColourY, y}).
		AddParameters(params...)
}

func NewColor(
	cmdID CommandID,
	level uint8,
	color color.Color,
	params ...Parameter,
) *Message {
	x, y := col2xy.Color2XY(color)
	return newColor(cmdID, level, x, y, params...)
}

func NewColorGroup(
	gid uint16,
	color color.Color,
	level uint8,
	params ...Parameter,
) *Message {
	return NewColor(DirectLevelGroup, level, color, params...).
		AddParameters(Parameter{Group, gid})
}

func NewColorDevice(
	addr string,
	color color.Color,
	level uint8,
	params ...Parameter,
) *Message {
	return NewColor(DirectLevelDevice, level, color, params...).
		AddParameters(Parameter{Address, addr})
}

func NewRGB(
	cmdID CommandID,
	level uint8,
	r, g, b byte,
	params ...Parameter,
) *Message {
	x, y := col2xy.RGB2XY(r, g, b)
	return newColor(cmdID, level, x, y, params...)
}

func NewRGBGroup(
	gid uint16,
	r, g, b byte,
	level uint8,
	params ...Parameter,
) *Message {
	return NewRGB(DirectLevelGroup, level, r, g, b, params...).
		AddParameters(Parameter{Group, gid})
}

func NewRGBDevice(
	addr string,
	r, g, b byte,
	level uint8,
	params ...Parameter,
) *Message {
	return NewRGB(DirectLevelDevice, level, r, g, b, params...).
		AddParameters(Parameter{Address, addr})
}
