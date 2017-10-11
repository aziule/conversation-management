package nlu

// Represents an entity, as understood by our NLU services
type Entity struct {
	name string
	value string
}

// Getters
func (entity *Entity) Name() string { return entity.name }
func (entity *Entity) Value() string { return entity.value }
