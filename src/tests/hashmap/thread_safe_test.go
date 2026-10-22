package hashmap

import (
	"fmt"
	"sync"
	"testing"
	h "treenity/src/server/hashmap"
)

func TestThreadSafe(t *testing.T) {
	myMap := h.CreateHashmap(1)

	var wg sync.WaitGroup
	workers := 50
	iterations := 100

	// Create a giant monoyy-hash hashmap
	for i := range workers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range iterations {
				key := fmt.Sprintf("key-%d-%d", workerID, j)
				myMap.Append(key, nil)
			}
		}(i)
	}
	wg.Wait()

	// Check init worked
	expectedLen := workers * iterations
	if myMap.Len() != expectedLen {
		t.Errorf("Wanted post appends length %d , got %d", expectedLen, myMap.Len())
	}

	// Reads simultaly to check collisions
	for i := range workers {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range iterations {
				key := fmt.Sprintf("key-%d-%d", workerID, j)

				// Read
				_, _ = myMap.Get(key)

				// Check data
				_ = myMap.Contains(key)

				// Delete
				if j%2 == 0 {
					_ = myMap.Delete(key)
				}
			}
		}(i)
	}
	wg.Wait()

	t.Logf("Concurrency tests finished. Final length : %d", myMap.Len())
}
