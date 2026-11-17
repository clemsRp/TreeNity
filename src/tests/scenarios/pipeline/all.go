package pipeline

import "slices"

var Scenarios = slices.Concat(Global, Message, EdgeCase)
