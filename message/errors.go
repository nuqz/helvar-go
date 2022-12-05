package message

type ErrorID uint8

const (
	EOK                     ErrorID = 0
	EInvalidGroupIndex      ErrorID = 1
	EInvalidCluster         ErrorID = 2
	EInvalidRouter          ErrorID = 3
	EInvalidSubnet          ErrorID = 4
	EInvalidDevice          ErrorID = 5
	EInvalidSubDevice       ErrorID = 6
	EInvalidBlock           ErrorID = 7
	EInvalidScene           ErrorID = 8
	EClusterDoesntExist     ErrorID = 9
	ERouterDoesntExist      ErrorID = 10
	EDeviceDoesntExist      ErrorID = 11
	EPropertyDoesntExist    ErrorID = 12
	EInvalidRawMessageSize  ErrorID = 13
	EInvalidMessagesType    ErrorID = 14
	EInvalidMessageCommand  ErrorID = 15
	EMissingASCIITerminator ErrorID = 16
	EMissingASCIIParameter  ErrorID = 17
	EIncompatibleVersion    ErrorID = 18
)

var ErrorsByID = map[ErrorID]error{}
