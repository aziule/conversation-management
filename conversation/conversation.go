package conversation

type Status string

const (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
)

// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Status       Status
	MessagesFlow *MessagesFlow
}

type Repository interface {
	FindLatestConversation(user *User) (*Conversation, error)
	SaveConversation(conversation *Conversation) error
	FindUser(userId string) (*User, error)
	InsertUser(user *User) error
}
