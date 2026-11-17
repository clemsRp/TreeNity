package pipeline

import (
	"encoding/binary"
	"fmt"
	"strings"
	"syscall"
	"time"
	pr "treenity/src/protocol"
	sce "treenity/src/tests/scenarios"
)

const (
	bigTimeout  = 15 * time.Second
	invalidName = "invalid name, must respect format " + pr.ValidRegex
)

func ok(args []string, out ...string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: args, WantOut: out}
}

func fails(code int, stderr string, args ...string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: args, WantCode: code, WantErr: stderr}
}

func create(topic string) sce.PipelineStep {
	return ok([]string{"create", topic}, "topic created")
}

func createAll(topics ...string) []sce.PipelineStep {
	var all []sce.PipelineStep
	for _, topic := range topics {
		all = append(all, create(topic))
	}
	return all
}

func list(topics ...string) sce.PipelineStep {
	if len(topics) == 0 {
		return ok([]string{"list"})
	}
	return ok([]string{"list"}, strings.Join(topics, ","))
}

func info(name, topic string, offset int, prefix string) sce.PipelineStep {
	want := fmt.Sprintf(`\{"client":"%s","topic":"%s","offset":%d,"prefix":"%s","ipc":"[^"]+"\}`,
		name, topic, offset, prefix)
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"info", name}, WantOut: []string{want}, WantRegexp: true}
}

func produce(topic, stdin string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", topic}, Stdin: stdin}
}

func produceRaw(topic, stdin string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", topic, "--raw"}, Stdin: stdin}
}

func produceBig(topic, stdin string) sce.PipelineStep {
	step := produce(topic, stdin)
	step.Timeout = bigTimeout
	return step
}

func produceFails(topic, stdin string, code int, stderr string, flags ...string) sce.PipelineStep {
	step := fails(code, stderr, append([]string{"produce", topic}, flags...)...)
	step.Stdin = stdin
	return step
}

func subscribe(proc, topic string, flags ...string) []sce.PipelineStep {
	args := append([]string{"subscribe", topic, proc}, flags...)
	return []sce.PipelineStep{
		{Action: sce.PipelineActionSubscribe, Proc: proc, Args: args},
		{Action: sce.PipelineActionExpect, Proc: proc, WantOut: []string{"subscribed to " + topic}},
	}
}

func subscribeMany(base, topic string, n int) []sce.PipelineStep {
	var all []sce.PipelineStep
	for i := range n {
		all = append(all, subscribe(fmt.Sprintf("%s%d", base, i), topic)...)
	}
	return all
}

func expect(proc string, lines ...string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionExpect, Proc: proc, WantOut: lines}
}

func expectBig(proc string, lines []string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionExpect, Proc: proc, WantOut: lines, Timeout: bigTimeout}
}

func expectMany(base string, n int, lines []string) []sce.PipelineStep {
	var all []sce.PipelineStep
	for i := range n {
		all = append(all, expectBig(fmt.Sprintf("%s%d", base, i), lines))
	}
	return all
}

func silent(proc string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionSilent, Proc: proc}
}

func signal(proc string, sig syscall.Signal) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionSignal, Proc: proc, Signal: sig}
}

func exits(proc string, code int, stderr string) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionWait, Proc: proc, WantCode: code, WantErr: stderr}
}

func interrupt(proc string) []sce.PipelineStep {
	return []sce.PipelineStep{signal(proc, syscall.SIGINT), exits(proc, 0, "")}
}

func kill(proc string) []sce.PipelineStep {
	return []sce.PipelineStep{signal(proc, syscall.SIGKILL), exits(proc, -1, "")}
}

func stopServer(sig syscall.Signal) sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionStopServer, Signal: sig}
}

func killServer() sce.PipelineStep {
	return sce.PipelineStep{Action: sce.PipelineActionStopServer, Signal: syscall.SIGKILL}
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
	return keyed("k", n)
}

func keyed(prefix string, n int) (string, []string) {
	var stdin strings.Builder
	lines := make([]string, n)
	for i := range n {
		lines[i] = fmt.Sprintf("%s%d:%s", prefix, i, strings.Repeat("v", 100))
		stdin.WriteString(lines[i] + "\n")
	}
	return stdin.String(), lines
}

func input(lines ...string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func record(key, body string) string {
	return length(len(key)) + key + length(len(body)) + body
}

func length(n int) string {
	return string(binary.LittleEndian.AppendUint32(nil, uint32(int32(n))))
}
