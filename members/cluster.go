package members

// Cluster represents a cluster as defined in HelvarNET protocol.
type Cluster struct {
	// ID is uint8, however HelvarNET protocol assumes it is within
	// [1..253] range.
	ID uint8
}
