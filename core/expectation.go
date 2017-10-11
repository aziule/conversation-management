package core

import "github.com/aziule/conversation-management/core/nlu"

//type expectation interface {
//	IsMet() bool
//}

// IntentExpectation represents an expectation with a single intent
type IntentExpectation struct {
	intent *nlu.Intent
}

// EntityExpectation represents an expectation where we have a single entity
type EntityExpectation struct {
	entity *nlu.Entity
}

/*
type condition interface {
	Validate() bool
	Left() *expectation
	Right() *expectation
}

type condition struct {
	left *expectation
	right *expectation
}

type andCondition struct {
	condition *condition
}

type orCondition struct {
	condition *condition
}
*/
