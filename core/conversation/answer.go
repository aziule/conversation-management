package conversation

import (
	"math/rand"
)

// AnswerPool represents a set of possible answers.
// Each story can have multiple AnswerGroup, identified by a unique
// name. Then, according to the kind of answer we want to send,
// we can choose one at random.
type AnswerPool struct {
	Name    string
	Answers []*Answer
}

// Answer is the main answer struct, containing the text to be sent.
type Answer struct {
	// Text contains the text to send.
	// It also contains placeholders.
	// @todo: define the placeholders, which ones are available, and their format.
	Text string
}

// RandomAnswer returns a random answer from a pool of answers.
// Returns nil if there is no answer available.
// @todo: test it
func (pool *AnswerPool) RandomAnswer() *Answer {
	nbAnswers := len(pool.Answers)

	if nbAnswers == 0 {
		return nil
	}

	return pool.Answers[rand.Intn(nbAnswers)]
}
