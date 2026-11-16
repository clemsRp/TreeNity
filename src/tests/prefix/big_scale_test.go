package prefix

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"

	pr "treenity/src/protocol"
	p "treenity/src/server/prefix"
)

const (
	bigSeed        = 42
	bigNbClients   = 10000
	bigNbKeys      = 3000
	noisePrefix    = "§noise§"
	maxSegmentsLen = 7
)

var (
	// Words used to build prefixes and keys
	bigWords = []string{
		"user", "users", "usr", "u", "admin", "adm", "sports", "sport", "news",
		"order", "orders", "ord", "payment", "pay", "log", "logs", "auth",
		"login", "logout", "create", "update", "delete", "a", "ab", "abc",
		"é", "éa", "日本", "日本語", "🙂", "🙂🙂", "ñandú", "0", "42", "x_y", "k-v",
	}

	// Separators between words (empty one glues words together)
	bigSeps = []string{".", ".", ".", "-", "_", ""}
)

func randPrefix(r *rand.Rand) string {
	// Create wildcard client
	if r.Intn(20) == 0 {
		return ""
	}

	// Get number of segments
	n := 1 + r.Intn(maxSegmentsLen)

	// Build prefix from random words
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString(bigSeps[r.Intn(len(bigSeps))])
		}
		sb.WriteString(bigWords[r.Intn(len(bigWords))])
	}
	s := sb.String()

	// Cut prefix at a random byte to get partial prefixes
	if r.Intn(4) == 0 {
		bytes := []byte(s)
		s = string(bytes[:r.Intn(len(bytes)+1)])
	}

	return s
}

func randKey(r *rand.Rand, clients []*pr.Client) string {
	// Half of the keys come from existing prefixes to get frequent matches
	switch r.Intn(4) {
	case 0:
		// Fully random key
		return randPrefix(r)

	case 1:
		// Exact existing prefix
		return clients[r.Intn(len(clients))].Prefix

	case 2:
		// Existing prefix with a suffix
		base := clients[r.Intn(len(clients))].Prefix
		return base + bigSeps[r.Intn(len(bigSeps))] + bigWords[r.Intn(len(bigWords))]

	default:
		// Truncated existing prefix
		bytes := []byte(clients[r.Intn(len(clients))].Prefix)
		return string(bytes[:r.Intn(len(bytes)+1)])
	}
}

func expectedIDs(alive []*pr.Client, key string) []string {
	var ids []string
	for _, c := range alive {
		if strings.HasPrefix(key, c.Prefix) {
			ids = append(ids, c.Id)
		}
	}
	sort.Strings(ids)

	return ids
}

func gotIDs(t *testing.T, got []*pr.Client, key string) []string {
	t.Helper()

	seen := make(map[string]bool, len(got))
	ids := make([]string, 0, len(got))

	for _, c := range got {
		// Handle duplicated client
		if seen[c.Id] {
			t.Errorf("Filter(%q) returned client %s more than once", key, c.Id)
			continue
		}
		seen[c.Id] = true
		ids = append(ids, c.Id)
	}
	sort.Strings(ids)

	return ids
}

func compareWithOracle(t *testing.T, tree *p.Tree, alive []*pr.Client, keys []string, label string) {
	t.Helper()

	for _, key := range keys {
		want := expectedIDs(alive, key)
		got := gotIDs(t, tree.Filter(key), key)

		// Check lengths
		if len(want) != len(got) {
			t.Fatalf("[%s] Filter(%q): got %d clients, want %d", label, key, len(got), len(want))
		}

		// Check same clients
		for i := range want {
			if want[i] != got[i] {
				t.Fatalf("[%s] Filter(%q): mismatch at %d: got %s, want %s", label, key, i, got[i], want[i])
			}
		}
	}
}

func aliveOf(clients []*pr.Client, isAlive map[string]bool) []*pr.Client {
	out := make([]*pr.Client, 0, len(clients))

	for _, c := range clients {
		if isAlive[c.Id] {
			out = append(out, c)
		}
	}

	return out
}

func TestBigScale(t *testing.T) {
	// Init random generator (fixed seed to reproduce failures)
	r := rand.New(rand.NewSource(bigSeed))
	t.Logf("seed=%d clients=%d keys=%d", bigSeed, bigNbClients, bigNbKeys)

	tree := p.CreateTree(0)

	// Insert all clients
	clients := make([]*pr.Client, 0, bigNbClients)
	isAlive := make(map[string]bool, bigNbClients)

	for i := 0; i < bigNbClients; i++ {
		prefix := randPrefix(r)
		c := &pr.Client{
			Id:             fmt.Sprintf("c%d", i),
			Prefix:         prefix,
			PreviousPrefix: prefix,
		}

		if err := tree.Insert(c); err != nil {
			t.Fatalf("Insert(%s, %q) failed: %v", c.Id, c.Prefix, err)
		}

		clients = append(clients, c)
		isAlive[c.Id] = true
	}

	// Check all clients are inserted
	if got := len(tree.GetAllClients()); got != bigNbClients {
		t.Fatalf("GetAllClients() = %d, want %d", got, bigNbClients)
	}

	// Generate keys (reused for each phase)
	keys := make([]string, 0, bigNbKeys+4)
	keys = append(keys, "", "zzzz", "\x00", "🙂🙂🙂🙂")
	for i := 0; i < bigNbKeys; i++ {
		keys = append(keys, randKey(r, clients))
	}

	compareWithOracle(t, tree, aliveOf(clients, isAlive), keys, "after insert")

	// Insert same clients twice (should not fail and keep the tree unchanged)
	for i := 0; i < 200; i++ {
		c := clients[r.Intn(len(clients))]
		if err := tree.Insert(c); err != nil {
			t.Fatalf("Insert(%s) twice should not fail (update)", c.Id)
		}
	}

	if got := len(tree.GetAllClients()); got != bigNbClients {
		t.Fatalf("after duplicate inserts GetAllClients() = %d, want %d", got, bigNbClients)
	}

	// Remove ~40% of clients
	for _, c := range clients {
		if r.Intn(100) < 40 {
			if err := tree.Remove(c); err != nil {
				t.Fatalf("Remove(%s, %q) failed: %v", c.Id, c.Prefix, err)
			}
			isAlive[c.Id] = false

			// Remove same client twice (must fail)
			if r.Intn(50) == 0 {
				if err := tree.Remove(c); err == nil {
					t.Fatalf("Remove(%s) twice should fail", c.Id)
				}
			}
		}
	}

	compareWithOracle(t, tree, aliveOf(clients, isAlive), keys, "after removals")

	// Change prefix of ~30% of alive clients
	for _, c := range aliveOf(clients, isAlive) {
		if r.Intn(100) < 30 {
			c.PreviousPrefix = c.Prefix
			c.Prefix = randPrefix(r)

			if err := tree.Insert(c); err != nil {
				t.Fatalf("Insert(update %s: %q -> %q) failed: %v", c.Id, c.PreviousPrefix, c.Prefix, err)
			}

			// Sync previous prefix like the server does on next subscribe
			c.PreviousPrefix = c.Prefix
		}
	}

	// Check no client was lost or duplicated by updates
	if got, want := len(tree.GetAllClients()), len(aliveOf(clients, isAlive)); got != want {
		t.Fatalf("after prefix updates GetAllClients() = %d, want %d", got, want)
	}

	compareWithOracle(t, tree, aliveOf(clients, isAlive), keys, "after prefix updates")

	// Reinsert half of removed clients with new prefixes
	for _, c := range clients {
		if !isAlive[c.Id] && r.Intn(2) == 0 {
			c.Prefix = randPrefix(r)
			c.PreviousPrefix = c.Prefix

			if err := tree.Insert(c); err != nil {
				t.Fatalf("re-Insert(%s, %q) failed: %v", c.Id, c.Prefix, err)
			}
			isAlive[c.Id] = true
		}
	}

	compareWithOracle(t, tree, aliveOf(clients, isAlive), keys, "after re-insertions")

	// Concurrent test
	stable := aliveOf(clients, isAlive)
	var wg sync.WaitGroup

	// Start readers
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			// Each reader has its own generator (rand.Rand is not thread safe)
			rr := rand.New(rand.NewSource(int64(1000 + worker)))

			for i := 0; i < 150; i++ {
				key := keys[rr.Intn(len(keys))]
				want := expectedIDs(stable, key)
				got := gotIDs(t, tree.Filter(key), key)

				if len(want) != len(got) {
					t.Errorf("[concurrent] Filter(%q): got %d, want %d", key, len(got), len(want))
					return
				}
			}
		}(w)
	}

	// Start writer
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 2000; i++ {
			prefix := fmt.Sprintf("%s%d", noisePrefix, i%50)
			c := &pr.Client{
				Id:             fmt.Sprintf("noise%d", i),
				Prefix:         prefix,
				PreviousPrefix: prefix,
			}

			if err := tree.Insert(c); err != nil {
				t.Errorf("[concurrent] noise Insert failed: %v", err)
				return
			}
			if err := tree.Remove(c); err != nil {
				t.Errorf("[concurrent] noise Remove failed: %v", err)
				return
			}
		}
	}()

	wg.Wait()

	compareWithOracle(t, tree, stable, keys, "after concurrency")

	// Remove all remaining clients
	for _, c := range stable {
		if err := tree.Remove(c); err != nil {
			t.Fatalf("final Remove(%s, %q) failed: %v", c.Id, c.Prefix, err)
		}
	}

	// Check tree is empty
	if got := len(tree.GetAllClients()); got != 0 {
		t.Fatalf("tree should be empty, still has %d clients", got)
	}

	// Check empty nodes were cleaned
	if len(tree.Children) != 0 {
		t.Fatalf("empty nodes were not cleaned: root still has %d children", len(tree.Children))
	}

	// Check filter returns nothing on empty tree
	for _, key := range keys {
		if got := tree.Filter(key); len(got) != 0 {
			t.Fatalf("Filter(%q) on empty tree returned %d clients", key, len(got))
		}
	}
}
