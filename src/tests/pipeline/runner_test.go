package pipeline

import (
	sce "treenity/src/tests/scenarios"
)

func (e *env) run(i int, step sce.PipelineStep) {
	// Handle timeout
	t := e.t
	timeout := step.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	// Exec corresponding action
	switch step.Action {
	case sce.PipelineActionRun:
		e.exec_run(t, step, timeout, i)

	case sce.PipelineActionSubscribe:
		e.exec_subscribe(t, step, timeout, i)

	case sce.PipelineActionExpect:
		e.exec_expect(t, step, timeout, i)

	case sce.PipelineActionSilent:
		e.exec_silent(t, step, timeout, i)

	case sce.PipelineActionSignal:
		e.exec_signal(t, step, timeout, i)

	case sce.PipelineActionWait:
		e.exec_wait(t, step, timeout, i)

	case sce.PipelineActionStopServer:
		e.exec_stop(t, step, timeout, i)

	default:
		t.Fatalf("step %d: unknown action %q", i, step.Action)
	}
}
