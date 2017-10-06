package nlu

// Represents an entity, as understood by our NLU services
type Entity struct {
	uuid string
	value string
}

// Getters
func (entity *Entity) Uuid() string { return entity.uuid }
func (entity *Entity) Value() string { return entity.value }
