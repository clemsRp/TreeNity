package pipeline

import (
	"bytes"
	"os/exec"
	"testing"
	"time"

	"treenity/src/tests/scenarios/pipeline"
)

const (
	defaultTimeout = 3 * time.Second
	silentDuration = 300 * time.Millisecond
)

type subscriber struct {
	cmd    *exec.Cmd
	lines  chan string
	stderr *bytes.Buffer
	exited chan struct{}
}

// State of one scenario
type env struct {
	t       *testing.T
	server  *exec.Cmd
	ipcPath string
	subs    map[string]*subscriber
}

func TestPipeline(t *testing.T) {
	for _, sc := range pipeline.Scenarios {
		t.Run(sc.Name, func(t *testing.T) {
			// Allow multiple test running at the same time
			t.Parallel()

			// Init server for each test
			server, ipcPath := startServer(t)
			e := &env{t: t, server: server, ipcPath: ipcPath, subs: map[string]*subscriber{}}
			t.Cleanup(e.killSubscribers)

			// Run test steps
			for i, step := range sc.Steps {
				e.run(i, step)
			}
		})
	}
}
