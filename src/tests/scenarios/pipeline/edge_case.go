package pipeline

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	sce "treenity/src/tests/scenarios"
)

// Big inputs, generated once
var (
	in3000, lines3000   = messages(3000)
	in6000, lines6000   = messages(6000)
	in1000, lines1000   = messages(1000)
	userIn, userLines   = keyed("user.", 3000)
	adminIn, adminLines = keyed("admin.", 3000)
)

// 50 topic names
var topics50 = func() []string {
	names := make([]string, 50)
	for i := range names {
		names[i] = fmt.Sprintf("t%d", i)
	}
	return names
}()

// Connect, receive one message and disconnect n times with the same name
func reconnectLoop(n int) []sce.PipelineStep {
	var all []sce.PipelineStep
	for i := range n {
		msg := fmt.Sprintf("k%d:v", i)
		all = append(all, steps(
			produce("t", msg+"\n"),
			subscribe("s1", "t"),
			expect("s1", msg),
			silent("s1"),
			interrupt("s1"),
		)...)
	}
	return all
}

var EdgeCase = []sce.PipelineScenario{
	// Big quantities
	{
		Name: "Replay 6000 messages",
		Steps: steps(
			create("t"),
			produceBig("t", in6000),
			subscribe("s1", "t"),
			expectBig("s1", lines6000),
			silent("s1"),
		),
	},
	{
		Name: "Live 6000 messages",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceBig("t", in6000),
			expectBig("s1", lines6000),
			silent("s1"),
		),
	},
	{
		Name: "Replay 3000 then live 3000",
		Steps: steps(
			create("t"),
			produceBig("t", in3000),
			subscribe("s1", "t"),
			produceBig("t", userIn),
			expectBig("s1", append(append([]string{}, lines3000...), userLines...)),
			silent("s1"),
		),
	},
	{
		Name: "Prefix filter on 6000 messages",
		Steps: steps(
			create("t"),
			produceBig("t", adminIn+userIn),
			subscribe("users", "t", "--prefix", "user."),
			subscribe("admins", "t", "--prefix", "admin."),
			expectBig("users", userLines),
			expectBig("admins", adminLines),
			silent("users"),
			silent("admins"),
			interrupt("users"),
			info("users", "t", 6000, "user."),
		),
	},
	{
		Name: "Live 1000 messages to 20 subscribers",
		Steps: steps(
			create("t"),
			subscribeMany("s", "t", 20),
			produceBig("t", in1000),
			expectMany("s", 20, lines1000),
		),
	},
	{
		Name: "Replay 3000 messages to 10 subscribers",
		Steps: steps(
			create("t"),
			produceBig("t", in3000),
			subscribeMany("s", "t", 10),
			expectMany("s", 10, lines3000),
		),
	},
	{
		Name: "50 topics",
		Steps: steps(
			createAll(topics50...),
			list(topics50...),
			fails(2, "'t49' topic already exists", "create", "t49"),
		),
	},
	{
		Name: "Same name reconnected 5 times",
		Steps: steps(
			create("t"),
			reconnectLoop(5),
			info("s1", "t", 5, ""),
		),
	},
	{
		Name: "Line over 64KB is refused",
		Steps: steps(
			create("t"),
			produceFails("t", "k:"+strings.Repeat("b", 70000)+"\n", 1, "token too long"),
		),
	},

	// Frozen subscriber (SIGSTOP)
	{
		Name: "Suspended subscriber does not block the server",
		Steps: steps(
			createAll("a", "b"),
			subscribe("stopped", "a"),
			subscribe("active", "a"),
			signal("stopped", syscall.SIGSTOP),
			sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", "a"}, Stdin: in3000, Timeout: 2 * time.Second},
			sce.PipelineStep{Action: sce.PipelineActionRun, Args: []string{"produce", "b"}, Stdin: "x:y\n", Timeout: time.Second},
			list("a", "b"),
			expectBig("active", lines3000),
		),
	},
	{
		Name: "Suspended then resumed subscriber gets everything",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			signal("s1", syscall.SIGSTOP),
			produceBig("t", in3000),
			signal("s1", syscall.SIGCONT),
			expectBig("s1", lines3000),
			silent("s1"),
		),
	},
	{
		Name: "Shutdown does not hang on a frozen subscriber",
		Steps: steps(
			create("t"),
			subscribe("frozen", "t"),
			subscribe("s1", "t"),
			signal("frozen", syscall.SIGSTOP),
			stopServer(syscall.SIGINT),
			exits("s1", 0, ""),
		),
	},

	// Killed subscriber (SIGKILL)
	{
		Name: "Killed subscriber can reconnect with the same name",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			kill("s1"),
			subscribe("s1", "t"),
			produce("t", "b:2\n"),
			expect("s1", "b:2"),
			silent("s1"),
		),
	},
	{
		Name: "Killed subscriber gets missed messages on reconnect",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			kill("s1"),
			produce("t", input("b:2", "c:3")),
			subscribe("s1", "t"),
			expect("s1", "b:2", "c:3"),
		),
	},
	{
		Name: "Killed subscriber does not disturb the others",
		Steps: steps(
			create("t"),
			subscribe("victim", "t"),
			subscribe("s1", "t"),
			kill("victim"),
			produceBig("t", in1000),
			expectBig("s1", lines1000),
			subscribe("s2", "t"),
			expectBig("s2", lines1000),
		),
	},
	{
		Name: "Subscriber killed during a big replay",
		Steps: steps(
			create("t"),
			produceBig("t", in6000),
			subscribe("victim", "t"),
			expect("victim", lines6000[:10]...),
			kill("victim"),
			subscribe("s1", "t"),
			expectBig("s1", lines6000),
			subscribe("victim", "t", "--offset", "0"),
			expectBig("victim", lines6000),
		),
	},
	{
		Name: "All subscribers killed",
		Steps: steps(
			create("t"),
			subscribeMany("s", "t", 5),
			kill("s0"), kill("s1"), kill("s2"), kill("s3"), kill("s4"),
			produce("t", "a:1\n"),
			subscribe("late", "t"),
			expect("late", "a:1"),
			stopServer(syscall.SIGINT),
			exits("late", 0, ""),
		),
	},

	// Server crash and closed server
	{
		Name: "Server killed: subscriber exits with code 3",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			killServer(),
			exits("s1", 3, "server down"),
		),
	},
	{
		Name: "Server killed: new commands fail",
		Steps: steps(
			create("t"),
			killServer(),
			fails(1, "failed to open fifo", "list"),
			fails(1, "failed to open fifo", "produce", "t"),
			fails(1, "failed to open fifo", "subscribe", "t", "s1"),
		),
	},
	{
		Name: "Server stopped: new commands fail",
		Steps: steps(
			create("t"),
			stopServer(syscall.SIGINT),
			fails(1, "failed to open fifo", "list"),
			fails(1, "failed to open fifo", "create", "u"),
		),
	},
	{
		Name: "Server shutdown during a big delivery",
		Steps: steps(
			create("t"),
			produceBig("t", in6000),
			subscribe("s1", "t"),
			expect("s1", lines6000[:10]...),
			stopServer(syscall.SIGINT),
			exits("s1", 0, ""),
		),
	},
}
