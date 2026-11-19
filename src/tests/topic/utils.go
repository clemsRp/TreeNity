package topic

import (
	"slices"
	"testing"
	t "treenity/src/server/topics"
	sce "treenity/src/tests/scenarios"
)

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func UtilTestAction(tt *testing.T, topics *t.Topics, saved map[string]*t.TopicDatas, step sce.TopicStep) {
	switch step.Action {
	case sce.TopicActionCreate:
		td := topics.Create(step.TopicName)
		if td == nil {
			tt.Fatalf("Create(%q) returned nil", step.TopicName)
		}
		if prev, ok := saved[step.TopicName]; ok && (prev == td) != step.WantSame {
			tt.Errorf("Create(%q) same pointer = %v, want %v", step.TopicName, prev == td, step.WantSame)
		}
		saved[step.TopicName] = td

	case sce.TopicActionDelete:
		err := topics.Delete(step.TopicName)
		if (err != nil) != step.WantErr {
			tt.Errorf("Delete(%q) error = %v, WantErr %v", step.TopicName, err, step.WantErr)
		}
		if err != nil && step.WantCode != 0 && err.Code != step.WantCode {
			tt.Errorf("Delete(%q) code = %d, want %d", step.TopicName, err.Code, step.WantCode)
		}

	case sce.TopicActionGet:
		_, exist := topics.Get(step.TopicName)
		if exist != step.WantExist {
			tt.Errorf("Get(%q) exist = %v, want %v", step.TopicName, exist, step.WantExist)
		}

	case sce.TopicActionList:
		got := topics.GetTopicNames()
		if !slices.Equal(got, step.WantNames) {
			tt.Errorf("GetTopicNames() = %v, want %v", got, step.WantNames)
		}

	case sce.TopicActionStopAll:
		topics.StopAll()

	case sce.TopicActionIsStopped:
		td := saved[step.TopicName]
		if td == nil {
			tt.Fatalf("IsStopped(%q): topic was never created in this scenario", step.TopicName)
		}
		if got := isClosed(td.Done); got != step.WantStopped {
			tt.Errorf("topic %q stopped = %v, want %v", step.TopicName, got, step.WantStopped)
		}
	}
}
