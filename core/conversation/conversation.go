package conversation

type Status string

const (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
)

type Conversation interface {
	Status() Status
	MessagesFlow() MessagesFlow
}
