package members

type Router struct {
	// ID is uint8, however HelvarNET protocol assumes it is within
	// [1..254] range.
	ID uint8
}
