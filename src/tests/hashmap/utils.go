package hashmap

import (
	"testing"
	pr "treenity/src/protocol"
	h "treenity/src/server/hashmap"
	sce "treenity/src/tests/scenarios"
)

func UtilTestAction(t *testing.T, myMap *h.Hashmap, step sce.HashmapStep) {
	switch step.Action {
	case sce.HashmapActionAppend:
		new_client := &pr.Client{Id: step.Value}
		myMap.Append(step.Key, new_client)

	case sce.HashmapActionGet:
		val, err := myMap.Get(step.Key)
		if (err != nil) != step.WantErr {
			t.Errorf("Get() error = %v, WantErr %v", err, step.WantErr)
		}
		if !step.WantErr && val.Id != step.Want {
			t.Errorf("Get() got = %v, want %v", val, step.Want)
		}
	}
}
