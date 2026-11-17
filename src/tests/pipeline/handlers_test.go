package pipeline

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
	sce "treenity/src/tests/scenarios"
)

func (e *env) exec_run(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, clientBin, append([]string{e.ipcPath}, step.Args...)...)
	cmd.Stdin = strings.NewReader(step.Stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	_ = cmd.Run()

	if ctx.Err() != nil {
		t.Fatalf("step %d: %v timed out after %v (server blocked?)", i, step.Args, timeout)
	}
	e.checkLines(i, step, splitLines(stdout.String()))
	e.checkStderr(i, step, stderr.String())
	e.checkCode(i, step, cmd.ProcessState.ExitCode(), stderr.String())
}

func (e *env) exec_subscribe(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	cmd := exec.Command(clientBin, append([]string{e.ipcPath}, step.Args...)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	sub := &subscriber{
		cmd:    cmd,
		lines:  make(chan string, 100000),
		stderr: &bytes.Buffer{},
		exited: make(chan struct{}),
	}
	cmd.Stderr = sub.stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Read stdout continuously, so the subscriber never blocks on it
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sub.lines <- scanner.Text()
		}
		close(sub.lines)
		_ = cmd.Wait()
		close(sub.exited)
	}()
	e.subs[step.Proc] = sub
}

func (e *env) exec_expect(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	sub := e.sub(i, step.Proc)
	deadline := time.After(timeout)
	var got []string
	for len(got) < len(step.WantOut) {
		select {
		case line, ok := <-sub.lines:
			if !ok {
				t.Fatalf("step %d: %s exited after %d/%d lines\nstderr: %s",
					i, step.Proc, len(got), len(step.WantOut), sub.stderr.String())
			}
			got = append(got, line)
		case <-deadline:
			t.Fatalf("step %d: %s received %d/%d lines in %v",
				i, step.Proc, len(got), len(step.WantOut), timeout)
		}
	}
	e.checkLines(i, step, got)
}

func (e *env) exec_silent(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	sub := e.sub(i, step.Proc)
	select {
	case line, ok := <-sub.lines:
		if ok {
			t.Fatalf("step %d: %s should receive nothing, got %q", i, step.Proc, line)
		}
	case <-time.After(silentDuration):
	}
}

func (e *env) exec_signal(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	if err := e.sub(i, step.Proc).cmd.Process.Signal(step.Signal); err != nil {
		t.Fatalf("step %d: %v", i, err)
	}
}

func (e *env) exec_wait(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	sub := e.sub(i, step.Proc)
	select {
	case <-sub.exited:
	case <-time.After(timeout):
		t.Fatalf("step %d: %s did not exit within %v", i, step.Proc, timeout)
	}
	e.checkStderr(i, step, sub.stderr.String())
	e.checkCode(i, step, sub.cmd.ProcessState.ExitCode(), sub.stderr.String())
	delete(e.subs, step.Proc)
}

func (e *env) exec_stop(t *testing.T, step sce.PipelineStep, timeout time.Duration, i int) {
	if err := e.server.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() { _ = e.server.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatalf("step %d: server did not stop", i)
	}
	if left, _ := filepath.Glob(e.ipcPath + "*"); len(left) > 0 {
		t.Fatalf("step %d: FIFOs left behind: %v", i, left)
	}
}
