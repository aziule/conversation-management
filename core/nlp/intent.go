package nlp

// Intent represents the underlying action of some text, as understood by the NLP service
type Intent struct {
	Name string `bson:"name"`
}

// NewIntent is the constructor method for Intent
func NewIntent(name string) *Intent {
	return &Intent{
		Name: name,
	}
}
