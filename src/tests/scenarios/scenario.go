package scenarios

import (
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
