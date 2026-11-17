package pipeline

import (
	"syscall"
	"time"

	sce "treenity/src/tests/scenarios"
)

var bigInput, bigLines = messages(3000)

var Scenarios = []sce.PipelineScenario{
	{
		Name: "Create and list topics",
		Steps: steps(
			create("alpha"),
			create("beta"),
			sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"list"}, WantOut: []string{"alpha,beta"}},
			fails(2, "'alpha' topic already exists", "create", "alpha"),
		),
	},
	{
		Name: "Unknown topic",
		Steps: steps(
			fails(2, "'nope' topic does not exist", "produce", "nope"),
			fails(2, "'nope' topic does not exist", "subscribe", "nope", "s1"),
		),
	},
	{
		Name: "Replay then live delivery",
		Steps: steps(
			create("t"),
			produce("t", "a:1\nb:2\n"),
			subscribe("s1", "t"),
			expect("s1", "a:1", "b:2"),
			produce("t", "c:3\n"),
			expect("s1", "c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Prefix filter",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "user."),
			produce("t", "user.a:1\nadmin.b:2\nuser.c:3\n"),
			expect("s1", "user.a:1", "user.c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Info on a subscriber",
		Steps: steps(
			create("t"),
			produce("t", "a:1\n"),
			subscribe("s1", "t"),
			expect("s1", "a:1"),
			sce.PipelineStep{
				Action:     sce.PipelineActionRun,
				Args:       []string{"info", "s1"},
				WantOut:    []string{`\{"client":"s1","topic":"t","offset":1,"prefix":"","ipc":"[^"]+"\}`},
				WantRegexp: true,
			},
		),
	},
	{
		Name: "Suspended subscriber does not block the server",
		Steps: steps(
			create("a"),
			create("b"),
			subscribe("stopped", "a"),
			subscribe("active", "a"),
			sce.PipelineStep{Action: sce.PipelineActionSignal, Proc: "stopped", Signal: syscall.SIGSTOP},
			sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", "a"}, Stdin: bigInput, Timeout: 2 * time.Second},
			sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", "b"}, Stdin: "x:y\n", Timeout: time.Second},
			sce.PipelineStep{Action: sce.PipelineActionExpect, Proc: "active", WantOut: bigLines},
		),
	},
	{
		Name: "Server shutdown",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			sce.PipelineStep{Action: sce.PipelineActionStopServer},
			sce.PipelineStep{Action: sce.PipelineActionWait, Proc: "s1", WantCode: 0},
		),
	},
}
