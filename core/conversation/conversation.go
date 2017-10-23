package conversation

type Status string

const (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
)

type Conversation interface {
	Channel() Channel
	Messages() []Message
	Status() Status
}
