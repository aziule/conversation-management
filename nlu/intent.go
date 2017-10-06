package nlu

// Represents an intent, as understood by our NLU services
type Intent struct {
	uuid string
}

// Getters
func (intent *Intent) Uuid() string { return intent.uuid }
