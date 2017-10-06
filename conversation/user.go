package conversation

// Our main User interface
type User interface {
	Uuid() string
	Name() string
}
