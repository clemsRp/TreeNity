package pipeline

import (
	"strings"
	"syscall"

	sce "treenity/src/tests/scenarios"
)

var Global = []sce.PipelineScenario{
	// create / list
	{
		Name: "List without topic prints nothing",
		Steps: steps(
			list(),
		),
	},
	{
		Name: "Create and list topics",
		Steps: steps(
			create("alpha"),
			create("beta"),
			list("alpha", "beta"),
		),
	},
	{
		Name: "List keeps the creation order",
		Steps: steps(
			createAll("zeta", "alpha", "mid"),
			list("zeta", "alpha", "mid"),
		),
	},
	{
		Name: "Create an existing topic fails",
		Steps: steps(
			create("alpha"),
			fails(2, "'alpha' topic already exists", "create", "alpha"),
			list("alpha"),
		),
	},
	{
		Name: "Topic names are case sensitive",
		Steps: steps(
			create("topic"),
			create("Topic"),
			list("topic", "Topic"),
		),
	},
	{
		Name: "Valid topic names",
		Steps: steps(
			create("a"),
			create("A-b_c.9"),
			create(strings.Repeat("x", 32)),
			list("a", "A-b_c.9", strings.Repeat("x", 32)),
		),
	},
	{
		Name: "Invalid topic names",
		Steps: steps(
			fails(1, "'bad name' "+invalidName, "create", "bad name"),
			fails(1, "'bad/name' "+invalidName, "create", "bad/name"),
			fails(1, "'été' "+invalidName, "create", "été"),
			fails(1, "'"+strings.Repeat("x", 33)+"' "+invalidName, "create", strings.Repeat("x", 33)),
			list(),
		),
	},

	// Unknown topic
	{
		Name: "Unknown topic",
		Steps: steps(
			fails(2, "'nope' topic does not exist", "produce", "nope"),
			fails(2, "'nope' topic does not exist", "subscribe", "nope", "s1"),
		),
	},
	{
		Name: "Invalid topic name on produce and subscribe",
		Steps: steps(
			fails(1, "'bad!' "+invalidName, "produce", "bad!"),
			fails(1, "'bad!' "+invalidName, "subscribe", "bad!", "s1"),
		),
	},
	{
		Name: "Invalid subscriber name",
		Steps: steps(
			create("t"),
			fails(1, "'bad name' "+invalidName, "subscribe", "t", "bad name"),
			fails(1, "'"+strings.Repeat("s", 33)+"' "+invalidName, "subscribe", "t", strings.Repeat("s", 33)),
		),
	},

	// info
	{
		Name: "Info on a new subscriber",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			info("s1", "t", 0, ""),
		),
	},
	{
		Name: "Info shows the prefix",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "user."),
			info("s1", "t", 0, "user."),
		),
	},
	{
		Name: "Info on a subscriber",
		Steps: steps(
			create("t"),
			produce("t", "a:1\n"),
			subscribe("s1", "t"),
			expect("s1", "a:1"),
			interrupt("s1"),
			info("s1", "t", 1, ""),
		),
	},
	{
		Name: "Info on an unknown subscriber",
		Steps: steps(
			fails(2, "'ghost' subscriber does not exist", "info", "ghost"),
		),
	},
	{
		Name: "Info with an invalid name",
		Steps: steps(
			fails(1, "'bad name' "+invalidName, "info", "bad name"),
		),
	},

	// Server stop
	{
		Name: "Server shutdown",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			stopServer(syscall.SIGINT),
			exits("s1", 0, ""),
		),
	},
	{
		Name: "Server shutdown with SIGTERM",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			stopServer(syscall.SIGTERM),
			exits("s1", 0, ""),
		),
	},
	{
		Name: "Server shutdown with several subscribers",
		Steps: steps(
			createAll("a", "b"),
			subscribe("s1", "a"),
			subscribe("s2", "a"),
			subscribe("s3", "b", "--prefix", "x"),
			stopServer(syscall.SIGINT),
			exits("s1", 0, ""),
			exits("s2", 0, ""),
			exits("s3", 0, ""),
		),
	},
	{
		Name: "Server shutdown without client",
		Steps: steps(
			createAll("a", "b"),
			produce("a", "k:v\n"),
			stopServer(syscall.SIGINT),
		),
	},
}
