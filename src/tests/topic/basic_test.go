package topic

import (
	"testing"
	tp "treenity/src/server/topics"
	sce "treenity/src/tests/scenarios/topic"
)

func TestBasic(t *testing.T) {
	tests := sce.ActionsScenarios

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			topics := tp.NewTopics()
			saved := make(map[string]*tp.TopicDatas)

			for _, step := range tt.Steps {
				UtilTestAction(t, topics, saved, step)
			}
		})
	}
}
