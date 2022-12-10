package members

// Group represents a group as defined in HelvarNET protocol.
type Group struct {
	// ID is uint16, however HelvarNET protocol assumes it is within
	// [1..16383] range.
	ID   uint16 `yaml:"id"`
	Name string `yaml:"name"`

	// Last scene is uint8, however HelvarNET protocol assumes it is within
	// [1..16] range.
	LastScene uint8 `yaml:"last_scene"`

	Devices []Device `yaml:"devices"`
}
