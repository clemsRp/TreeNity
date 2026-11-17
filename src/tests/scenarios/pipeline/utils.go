package pipeline

import (
	"fmt"
	"strings"
	sce "treenity/src/tests/scenarios"
)

func create(topic string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"create", topic}, WantOut: []string{"topic created"}}
}

func produce(topic, stdin string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", topic}, Stdin: stdin}
}

func subscribe(proc, topic string, flags ...string) []sce.PipelineStep {
	args := append([]string{"subscribe", topic, proc}, flags...)
	return []sce.PipelineStep{
		{Action: sce.PipelineActionSubscribe, Proc: proc, Args: args},
		{Action: sce.PipelineActionExpect, Proc: proc, WantOut: []string{"subscribed to " + topic}},
	}
}

func expect(proc string, lines ...string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionExpect, Proc: proc, WantOut: lines}
}

func silent(proc string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionSilent, Proc: proc}
}

func fails(code int, stderr string, args ...string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: args, WantCode: code, WantErr: stderr}
}

func steps(items ...any) []sce.PipelineStep {
	var all []sce.PipelineStep
	for _, item := range items {
		switch s := item.(type) {
		case sce.PipelineStep:
			all = append(all, s)
		case []sce.PipelineStep:
			all = append(all, s...)
		default:
			panic(fmt.Sprintf("steps: unexpected type %T", item))
		}
	}
	return all
}

func messages(n int) (string, []string) {
	var stdin strings.Builder
	lines := make([]string, n)
	for i := range n {
		lines[i] = fmt.Sprintf("k%d:%s", i, strings.Repeat("v", 100))
		stdin.WriteString(lines[i] + "\n")
	}
	return stdin.String(), lines
}
