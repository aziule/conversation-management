package conversation

type Status string

const (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
)

type Conversation struct {
	Status       Status
	MessagesFlow *MessagesFlow
}

type ConversationReader interface {
	FindLatest(*User) (*Conversation, error)
}
