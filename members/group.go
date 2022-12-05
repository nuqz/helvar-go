package members

// Group represents a group as defined in HelvarNET protocol.
type Group struct {
	// ID is uint16, however HelvarNET protocol assumes it is within
	// [1..16383] range.
	ID   uint16
	Name string

	// Last scene is uint8, however HelvarNET protocol assumes it is within
	// [1..16] range.
	LastScene uint8

	Devices []Device
}
