package scenarios

import (
	"syscall"
	"time"
	pr "treenity/src/protocol"
)

// Common scenario type
type BasicScenario struct {
	Name    string
	Args    string
	WantErr bool
}

// Hashmap scenarios
type HashmapAction string

const (
	HashmapActionAppend   HashmapAction = "append"
	HashmapActionGet      HashmapAction = "get"
	HashmapActionDelete   HashmapAction = "delete"
	HashmapActionLen      HashmapAction = "len"
	HashmapActionContains HashmapAction = "contains"
)

type HashmapStep struct {
	Action   HashmapAction
	Key      string
	Value    string
	Want     any
	WantErr  bool
	WantBool bool
	WantInt  int
}

type HashmapScenario struct {
	Name  string
	Steps []HashmapStep
}

// Prefix Tree scenarios
type PrefixAction string

const (
	PrefixActionInsert PrefixAction = "insert"
	PrefixActionRemove PrefixAction = "remove"
	PrefixActionFilter PrefixAction = "filter"
)

type PrefixStep struct {
	Action       PrefixAction
	Key          string
	Client       *pr.Client
	WantClients  []*pr.Client
	WantErr      bool
	CustomPrefix string
}

type PrefixScenario struct {
	Name  string
	Steps []PrefixStep
}

// IPC scenarios
type IpcAction string

const (
	IpcActionCreateReader IpcAction = "create_reader"
	IpcActionCreateWriter IpcAction = "create_writer"
	IpcActionRead         IpcAction = "read"
	IpcActionWrite        IpcAction = "write"
)

type IpcStep struct {
	Action  IpcAction
	Pid     int
	Text    string
	WantErr bool
}

type IpcScenario struct {
	Name  string
	Steps []IpcStep
}

// Pipeline scenarios
type PipelineAction string

const (
	PipelineActionRun        PipelineAction = "run"
	PipelineActionSubscribe  PipelineAction = "subscribe"
	PipelineActionExpect     PipelineAction = "expect"
	PipelineActionSilent     PipelineAction = "silent"
	PipelineActionSignal     PipelineAction = "signal"
	PipelineActionWait       PipelineAction = "wait"
	PipelineActionStopServer PipelineAction = "stop_server"
)

type PipelineStep struct {
	Action PipelineAction

	Args  []string
	Stdin string

	Proc   string
	Signal syscall.Signal

	WantOut    []string
	WantRegexp bool
	WantErr    string
	WantCode   int

	Timeout time.Duration
}

type PipelineScenario struct {
	Name  string
	Steps []PipelineStep
}
