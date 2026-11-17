package pipeline

import (
	"regexp"
	"strings"
	sce "treenity/src/tests/scenarios"
)

func (e *env) sub(i int, name string) *subscriber {
	sub, ok := e.subs[name]
	if !ok {
		e.t.Fatalf("step %d: no subscriber named %q", i, name)
	}
	return sub
}

func (e *env) checkLines(i int, step sce.PipelineStep, got []string) {
	if len(got) != len(step.WantOut) {
		e.t.Fatalf("step %d: got %d lines, want %d\ngot:  %q\nwant: %q",
			i, len(got), len(step.WantOut), got, step.WantOut)
	}
	for j, want := range step.WantOut {
		ok := got[j] == want
		if step.WantRegexp {
			ok = regexp.MustCompile("^(?:" + want + ")$").MatchString(got[j])
		}
		if !ok {
			e.t.Fatalf("step %d: line %d is %q, want %q", i, j, got[j], want)
		}
	}
}

func (e *env) checkStderr(i int, step sce.PipelineStep, stderr string) {
	stderr = strings.TrimSpace(stderr)
	if step.WantErr == "" && stderr != "" {
		e.t.Fatalf("step %d: unexpected stderr: %s", i, stderr)
	}
	if !strings.Contains(stderr, step.WantErr) {
		e.t.Fatalf("step %d: stderr is %q, want %q", i, stderr, step.WantErr)
	}
}

func (e *env) checkCode(i int, step sce.PipelineStep, code int, stderr string) {
	if code != step.WantCode {
		e.t.Fatalf("step %d: exit code %d, want %d\nstderr: %s", i, code, step.WantCode, stderr)
	}
}

func (e *env) killSubscribers() {
	for _, sub := range e.subs {
		_ = sub.cmd.Process.Kill()
		<-sub.exited
	}
}

func splitLines(s string) []string {
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
