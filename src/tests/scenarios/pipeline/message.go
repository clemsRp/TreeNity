package pipeline

import (
	"strings"
	"syscall"

	sce "treenity/src/tests/scenarios"
)

var Message = []sce.PipelineScenario{
	// Delivery
	{
		Name: "Replay then live delivery",
		Steps: steps(
			create("t"),
			produce("t", "a:1\nb:2\n"),
			subscribe("s1", "t"),
			expect("s1", "a:1", "b:2"),
			produce("t", "c:3\n"),
			expect("s1", "c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Live delivery only",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			silent("s1"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			produce("t", "b:2\nc:3\n"),
			expect("s1", "b:2", "c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Order is kept",
		Steps: steps(
			create("t"),
			produce("t", input("k0:a", "k1:b", "k2:c", "k3:d", "k4:e")),
			produce("t", input("k5:f", "k6:g")),
			subscribe("s1", "t"),
			expect("s1", "k0:a", "k1:b", "k2:c", "k3:d", "k4:e", "k5:f", "k6:g"),
		),
	},
	{
		Name: "Every subscriber receives every message",
		Steps: steps(
			create("t"),
			produce("t", "a:1\n"),
			subscribe("s1", "t"),
			subscribe("s2", "t"),
			subscribe("s3", "t"),
			produce("t", "b:2\n"),
			expect("s1", "a:1", "b:2"),
			expect("s2", "a:1", "b:2"),
			expect("s3", "a:1", "b:2"),
		),
	},
	{
		Name: "Topics are isolated",
		Steps: steps(
			createAll("a", "b"),
			subscribe("sa", "a"),
			subscribe("sb", "b"),
			produce("a", "x:1\n"),
			expect("sa", "x:1"),
			silent("sb"),
			produce("b", "y:2\n"),
			expect("sb", "y:2"),
			silent("sa"),
		),
	},
	{
		Name: "Same key several times",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", input("k:1", "k:2", "k:1")),
			expect("s1", "k:1", "k:2", "k:1"),
		),
	},

	// Text format
	{
		Name: "Body containing ':'",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", input("url:http://host:8080/path", "time:12:30:00")),
			expect("s1", "url:http://host:8080/path", "time:12:30:00"),
		),
	},
	{
		Name: "Empty key or empty body",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", input("key:", ":body", ":")),
			expect("s1", "key:", ":body", ":"),
		),
	},
	{
		Name: "Spaces and unicode",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", input("user 1:hello world", "clé:héllo wörld 🚀", "tab:a\tb")),
			expect("s1", "user 1:hello world", "clé:héllo wörld 🚀", "tab:a\tb"),
		),
	},
	{
		Name: "Last line without newline",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\nb:2"),
			expect("s1", "a:1", "b:2"),
		),
	},
	{
		Name: "Empty input sends nothing",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", ""),
			silent("s1"),
		),
	},
	{
		Name: "Line without ':' stops the producer",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceFails("t", input("a:1", "no separator", "c:3"), 1, "invalid format, must respect format 'key:body'"),
			expect("s1", "a:1"),
			silent("s1"),
		),
	},
	{
		Name: "Empty line stops the producer",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceFails("t", input("a:1", "", "c:3"), 1, "invalid format"),
			expect("s1", "a:1"),
			silent("s1"),
		),
	},
	{
		Name: "Message of exactly 1024 bytes",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "k:"+strings.Repeat("b", 1023)+"\n"),
			expect("s1", "k:"+strings.Repeat("b", 1023)),
		),
	},
	{
		Name: "Message over 1024 bytes is refused",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceFails("t", input("a:1", "k:"+strings.Repeat("b", 1024), "c:3"), 1, "message length exceeds 1024 bytes"),
			expect("s1", "a:1"),
			silent("s1"),
		),
	},

	// Raw format (--raw)
	{
		Name: "Raw produce",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceRaw("t", record("a", "1")+record("key", "hello world")+record("u", "x:y")),
			expect("s1", "a:1", "key:hello world", "u:x:y"),
		),
	},
	{
		Name: "Raw empty input sends nothing",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceRaw("t", ""),
			silent("s1"),
		),
	},
	{
		Name: "Raw partial record",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceFails("t", record("a", "1")+length(5)+"ab", 1, "partial record at EOF", "--raw"),
			expect("s1", "a:1"),
			silent("s1"),
		),
	},
	{
		Name: "Raw truncated length",
		Steps: steps(
			create("t"),
			produceFails("t", record("a", "1")+"\x01\x00", 1, "partial record at EOF", "--raw"),
		),
	},
	{
		Name: "Raw negative length",
		Steps: steps(
			create("t"),
			produceFails("t", length(-1), 1, "partial record at EOF", "--raw"),
		),
	},
	{
		Name: "Raw message over 1024 bytes",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produceFails("t", record("a", "1")+record(strings.Repeat("k", 1000), strings.Repeat("b", 25)), 1, "message length exceeds 1024 bytes", "--raw"),
			produceFails("t", length(1025), 1, "message length exceeds 1024 bytes", "--raw"),
			expect("s1", "a:1"),
			silent("s1"),
		),
	},

	// Prefix
	{
		Name: "Prefix filter",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "user."),
			produce("t", "user.a:1\nadmin.b:2\nuser.c:3\n"),
			expect("s1", "user.a:1", "user.c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Prefix filter on replay then live",
		Steps: steps(
			create("t"),
			produce("t", input("user.a:1", "admin.b:2")),
			subscribe("s1", "t", "--prefix", "user."),
			expect("s1", "user.a:1"),
			produce("t", input("admin.c:3", "user.d:4")),
			expect("s1", "user.d:4"),
			silent("s1"),
		),
	},
	{
		Name: "Overlapping prefixes",
		Steps: steps(
			create("t"),
			subscribe("u", "t", "--prefix", "u"),
			subscribe("us", "t", "--prefix", "us"),
			subscribe("user", "t", "--prefix", "user."),
			subscribe("all", "t"),
			produce("t", input("u:1", "us:2", "user.x:3", "user:4", "other:5")),
			expect("u", "u:1", "us:2", "user.x:3", "user:4"),
			expect("us", "us:2", "user.x:3", "user:4"),
			expect("user", "user.x:3"),
			expect("all", "u:1", "us:2", "user.x:3", "user:4", "other:5"),
			silent("u"),
			silent("us"),
			silent("user"),
		),
	},
	{
		Name: "Prefix equal to the whole key",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "user"),
			produce("t", input("use:1", "user:2", "users:3")),
			expect("s1", "user:2", "users:3"),
		),
	},
	{
		Name: "Empty prefix receives everything",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", ""),
			produce("t", input("a:1", ":2")),
			expect("s1", "a:1", ":2"),
		),
	},
	{
		Name: "Prefix is case sensitive",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "User"),
			produce("t", input("user:1", "User:2")),
			expect("s1", "User:2"),
			silent("s1"),
		),
	},

	// Offset
	{
		Name: "Start at an offset",
		Steps: steps(
			create("t"),
			produce("t", input("a:0", "b:1", "c:2", "d:3")),
			subscribe("s1", "t", "--offset", "2"),
			expect("s1", "c:2", "d:3"),
			silent("s1"),
		),
	},
	{
		Name: "Offset at the end only gets live messages",
		Steps: steps(
			create("t"),
			produce("t", input("a:0", "b:1")),
			subscribe("s1", "t", "--offset", "2"),
			silent("s1"),
			produce("t", "c:2\n"),
			expect("s1", "c:2"),
		),
	},
	{
		Name: "Offset after the end waits for that offset",
		Steps: steps(
			create("t"),
			produce("t", input("a:0", "b:1")),
			subscribe("s1", "t", "--offset", "4"),
			produce("t", input("c:2", "d:3")),
			silent("s1"),
			produce("t", "e:4\n"),
			expect("s1", "e:4"),
		),
	},
	{
		Name: "Offset and prefix together",
		Steps: steps(
			create("t"),
			produce("t", input("user.a:0", "admin:1", "user.b:2", "user.c:3")),
			subscribe("s1", "t", "--offset", "1", "--prefix", "user."),
			expect("s1", "user.b:2", "user.c:3"),
		),
	},
	{
		Name: "Info offset follows the acks",
		Steps: steps(
			create("t"),
			produce("t", input("a:0", "b:1", "c:2")),
			subscribe("s1", "t"),
			expect("s1", "a:0", "b:1", "c:2"),
			interrupt("s1"),
			info("s1", "t", 3, ""),
		),
	},
	{
		Name: "Info offset with a prefix stops at the last delivered message",
		Steps: steps(
			create("t"),
			produce("t", input("user.a:0", "admin.b:1", "user.c:2", "admin.d:3")),
			subscribe("s1", "t", "--prefix", "user."),
			expect("s1", "user.a:0", "user.c:2"),
			interrupt("s1"),
			info("s1", "t", 3, "user."),
		),
	},

	// Reconnection
	{
		Name: "Reconnect gets the messages produced while offline",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			interrupt("s1"),
			produce("t", input("b:2", "c:3")),
			subscribe("s1", "t"),
			expect("s1", "b:2", "c:3"),
			silent("s1"),
		),
	},
	{
		Name: "Reconnect with SIGTERM",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
			signal("s1", syscall.SIGTERM),
			exits("s1", 0, ""),
			produce("t", "b:2\n"),
			subscribe("s1", "t"),
			expect("s1", "b:2"),
		),
	},
	{
		Name: "Reconnect with --offset 0 replays everything",
		Steps: steps(
			create("t"),
			produce("t", input("a:1", "b:2")),
			subscribe("s1", "t"),
			expect("s1", "a:1", "b:2"),
			interrupt("s1"),
			subscribe("s1", "t", "--offset", "0"),
			expect("s1", "a:1", "b:2"),
			silent("s1"),
		),
	},
	{
		Name: "Reconnect with a new prefix",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "a."),
			produce("t", input("a.1:x", "b.1:y")),
			expect("s1", "a.1:x"),
			interrupt("s1"),
			info("s1", "t", 1, "a."),
			produce("t", input("a.2:x", "b.2:y")),
			subscribe("s1", "t", "--prefix", "b."),
			expect("s1", "b.1:y", "b.2:y"),
			silent("s1"),
			info("s1", "t", 4, "b."),
		),
	},
	{
		Name: "Reconnect keeps the previous prefix",
		Steps: steps(
			create("t"),
			subscribe("s1", "t", "--prefix", "a."),
			produce("t", "a.1:x\n"),
			expect("s1", "a.1:x"),
			interrupt("s1"),
			produce("t", input("b.1:y", "a.2:x")),
			subscribe("s1", "t"),
			expect("s1", "a.2:x"),
			silent("s1"),
		),
	},
	{
		Name: "Reconnect on another topic starts at 0",
		Steps: steps(
			createAll("a", "b"),
			produce("b", input("b:0", "b:1")),
			subscribe("s1", "a"),
			produce("a", input("a:0", "a:1", "a:2")),
			expect("s1", "a:0", "a:1", "a:2"),
			interrupt("s1"),
			subscribe("s1", "b"),
			expect("s1", "b:0", "b:1"),
			produce("a", "a:3\n"),
			silent("s1"),
			info("s1", "b", 2, ""),
		),
	},
	{
		Name: "Same subscriber twice at the same time is refused",
		Steps: steps(
			create("t"),
			subscribe("s1", "t"),
			fails(2, "'s1' subscriber is already active", "subscribe", "t", "s1"),
			produce("t", "a:1\n"),
			expect("s1", "a:1"),
		),
	},
	{
		Name: "Subscriber names are independent of topics",
		Steps: steps(
			createAll("a", "b"),
			subscribe("s1", "a"),
			fails(2, "'s1' subscriber is already active", "subscribe", "b", "s1"),
			subscribe("s2", "b"),
			produce("b", "x:1\n"),
			expect("s2", "x:1"),
			silent("s1"),
		),
	},
}
