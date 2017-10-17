package nlp

// Intent represents the underlying action of some text, as understood by the NLP service
type Intent struct {
	name string
}

// NewIntent is the constructor method for Intent
func NewIntent(name string) *Intent {
	return &Intent{
		name: name,
	}
}
