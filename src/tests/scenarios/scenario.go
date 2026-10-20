package scenarios

type BasicScenario struct {
	Name    string
	Args    string
	WantErr bool
}

type HashmapAction string

const (
	ActionAppend   HashmapAction = "append"
	ActionGet      HashmapAction = "get"
	ActionDelete   HashmapAction = "delete"
	ActionLen      HashmapAction = "len"
	ActionContains HashmapAction = "contains"
)

type HashmapStep struct {
	Action   HashmapAction
	Key      string
	Value    any
	Want     any
	WantErr  bool
	WantBool bool
	WantInt  int
}

type HashmapScenario struct {
	Name  string
	Steps []HashmapStep
}
