// Package memory implements objects working with data stored in memory.
package memory

import (
	"github.com/aziule/conversation-management/core/conversation"
	log "github.com/sirupsen/logrus"
)

// stories will store our stories
var stories []*conversation.Story

// inMemoryStoryRepository is the in memory implementation of a StoryRepository
type inMemoryStoryRepository struct{}

// NewStoryRepository instanciates a new in memory step repository
func NewStoryRepository() *inMemoryStoryRepository {
	return &inMemoryStoryRepository{}
}

// FindAll returns the full list of stories with the populated steps.
// In this repository, we use a static version of stories and steps.
func (r *inMemoryStoryRepository) FindAll() ([]*conversation.Story, error) {
	if stories != nil {
		log.WithField("stories", stories).Debug("Returning already fetched stories")
		return stories, nil
	}

	story := conversation.NewStory("Book a table", nil)
	step1 := conversation.NewStep(
		"book_table_entrypoint",
		"book_table",
		nil,
		nil,
	)

	step11 := conversation.NewStep(
		"book_table_get_nb_persons",
		"",
		[]string{"nb_persons"},
		nil,
	)

	step12 := conversation.NewStep(
		"book_table_get_time",
		"",
		[]string{"booking_date"},
		nil,
	)

	step1.AddNextStep(step11)
	step1.AddNextStep(step12)

	story.AddStartingStep(step1)

	stories = append(stories, story)

	log.WithField("stories", stories).Debug("Returning stories for the first time")

	return stories, nil
}
