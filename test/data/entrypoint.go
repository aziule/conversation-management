package data

import "github.com/aziule/conversation-management/conversation"

// Get a dummy entrypoint for tests. Schema:
// 1 -> 1.1 -> 1.1.1
// 1 -> 1.2 -> 2.1
// 2 -> 2.1
func GetDummyEntrypoint() *conversation.Entrypoint {
	var stories []*conversation.Story
	var startingSteps []*conversation.Step

	step21 := conversation.NewStep("Step 2.1", nil)
	step2 := conversation.NewStep("Step 2", []*conversation.Step{step21})
	step12 := conversation.NewStep("Step 1.2", []*conversation.Step{step21})
	step111 := conversation.NewStep("Step 1.1.1", nil)
	step11 := conversation.NewStep("Step 1.1", []*conversation.Step{step111})
	step1 := conversation.NewStep("Step 1", []*conversation.Step{step11, step12})

	startingSteps = append(startingSteps, step1)
	startingSteps = append(startingSteps, step2)

	story1 := conversation.NewStory("Test story", startingSteps)
	stories = append(stories, story1)

	entrypoint := conversation.NewEntryPoint(stories)

	return entrypoint
}