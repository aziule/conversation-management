package nlu

// Intent represents the data of a text, as understood by the NLU service
type Entity struct {
	name string
	value string
}

// NewEntity is the constructor method for Entity
func NewEntity(name, value string) *Entity {
	return &Entity{
		name: name,
		value: value,
	}
}

// Getters
func (entity *Entity) Name() string { return entity.name }
func (entity *Entity) Value() string { return entity.value }
