package hashmap

import (
	"testing"
	h "treenity/src/server/hashmap"
	sce "treenity/src/tests/scenarios"
)

func UtilTestAction(t *testing.T, myMap *h.Hashmap, step sce.HashmapStep) {
	switch step.Action {
	case sce.HashmapActionAppend:
		myMap.Append(step.Key, step.Value)

	case sce.HashmapActionGet:
		val, err := myMap.Get(step.Key)
		if (err != nil) != step.WantErr {
			t.Errorf("Get() error = %v, WantErr %v", err, step.WantErr)
		}
		if !step.WantErr && val != step.Want {
			t.Errorf("Get() got = %v, want %v", val, step.Want)
		}

	case sce.HashmapActionDelete:
		err := myMap.Delete(step.Key)
		if (err != nil) != step.WantErr {
			t.Errorf("Delete() error = %v, WantErr %v", err, step.WantErr)
		}

	case sce.HashmapActionContains:
		exists := myMap.Contains(step.Key)
		if exists != step.WantBool {
			t.Errorf("Contains() for key '%s' got = %v, want %v", step.Key, exists, step.WantBool)
		}

	case sce.HashmapActionLen:
		length := myMap.Len()
		if length != step.WantInt {
			t.Errorf("Length() got = %d, want %d", length, step.WantInt)
		}
	}
}
