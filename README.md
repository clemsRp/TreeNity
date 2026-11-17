*This project has been created as part of the 42 curriculum by crappo, lgoderne, opernod.*

# treenity

## Description

treenity is a lightweight, concurrent message queue system implemented in Go, inspired by the publish-subscribe messaging architecture of modern event-streaming platforms.

The system is built around a centralized server process and a command-line client interface. It enables asynchronous, decoupled communication between multiple producers and consumers through topics. Core capabilities of treenity include:

- Topic Management: Dynamic creation and listing of isolated message topics.
- Publish-Subscribe Messaging: Multiple producers can push messages to a topic while multiple consumers independently read and process messages.
- Prefix Matching: Consumers can subscribe to a topic with an optional key prefix filter, receiving only messages whose keys start with that prefix.
- Offset Tracking: Consumers track their position within a topic through committed offsets, allowing subscribers to resume reading from where they left off or replay past messages.
- Dual Message Formats: Support for both human-readable text format (`key:body`) and high-throughput raw binary format (`[keysize:int32][key:bytes][valuesize:int32][value:bytes]`).
- Robust Concurrency: Thread-safe data structures and isolated worker routines ensuring high throughput without blocking the central server or independent consumers.

## Instructions

### Compilation and Build Rules

The project includes a Makefile at the root directory managing compilation, code hygiene, and testing:

- Compile both server and client binaries:
```bash
make
# or
make all
```

- Compile only the server executable:
```bash
make server
```

- Compile only the client executable:
```bash
make client
```

- Clean build artifacts:
```bash
make clean
```

- Rebuild the entire project:
```bash
make re
```

- Format and lint source code:
```bash
make format
make lint
```

### Execution

#### 1. Starting the Server

The server acts as the central coordinator. It must be started first:

```bash
./server
```

Upon starting, the server prints its unique IPC entrypoint path to standard output:

```text
/tmp/operator.server.<PID>
```

The server listens for incoming client requests, manages topics and client registries, and gracefully shuts down upon receiving `SIGINT` or `SIGTERM`.

#### 2. Using the Client

All client operations interact with the server via the IPC entrypoint identifier returned by the server.

- Create a Topic:
```bash
./client <ipc_identifier> create <topic_name>
```
Output:
```text
topic created
```

- List Available Topics:
```bash
./client <ipc_identifier> list
```
Output:
```text
topic1,topic2,topic3
```

- Produce Messages (Text Mode):
Producers read lines from standard input in `key:body` format:
```bash
./client <ipc_identifier> produce <topic_name>
```
Example interactive input:
```text
user.signup:{"id": 42, "name": "alice"}
user.login:{"id": 42}
order.created:{"order_id": 1001, "amount": 59.99}
```

- Produce Messages (Raw Binary Mode):
Using the `--raw` flag, the producer accepts binary records encoded with little-endian 32-bit integers:
`[keysize:int32][key:bytes][valuesize:int32][value:bytes]`
```bash
./client <ipc_identifier> produce <topic_name> --raw
```

- Subscribe to Messages:
Consumers subscribe to a topic using a unique subscriber name:
```bash
./client <ipc_identifier> subscribe <topic_name> <subscriber_name>
```
Options:
- `--prefix <prefix>`: Filter messages by key prefix.
- `--offset <offset>`: Start reading from a specific offset (defaults to 0 for new consumers or the last committed offset for returning consumers).
- `--raw`: Receive messages in binary format `[offset:int32][keysize:int32][key:bytes][valuesize:int32][value:bytes]`.

Example:
```bash
./client <ipc_identifier> subscribe user_events client_alpha --prefix user. --offset 0
```
Output:
```text
subscribed to user_events
user.signup:{"id": 42, "name": "alice"}
user.login:{"id": 42}
```

- Query Client Metadata (Info):
Inspect the current state and offset of a registered subscriber:
```bash
./client <ipc_identifier> info <subscriber_name>
```
Output (JSON):
```json
{"client":"client_alpha","topic":"user_events","offset":2,"prefix":"user.","ipc":"/tmp/operator.server.1234.5678"}
```

### Exit Codes

- `0`: Success or clean graceful shutdown.
- `1`: General error (invalid arguments, malformed message format, partial binary record at EOF, invalid identifier name).
- `2`: Topic error (topic does not exist when producing/subscribing, or topic already exists when creating).
- `3`: IPC Comunication error (FileOpen error, invlid filepath, invalid rights...).

## Architecture

treenity is designed around a multi-process architecture where a centralized server coordinates independent client processes via dedicated communication channels.

### Threading and Goroutine Model per Topic

To ensure high performance and prevent bottlenecks, concurrency is managed using Go goroutines structured per topic and per subscriber:

1. Central Request Listener:
The main server routine runs an asynchronous reader loop (`ReadChannel`) on the main entrypoint FIFO. Incoming requests (topic creation, topic listing, subscriber registration, metadata queries) are received and dispatched concurrently without blocking the main event loop.

2. Dedicated Topic Distributor Routine:
When a topic is created (`HandleCmdCreate`), the server spawns a dedicated goroutine:
```go
go new_topic.DistributeMessages(clients)
```
This routine owns message ingestion for that specific topic. Incoming messages arrive over an unbuffered channel (`td.Topic.NewMsgChan`). This ensures that message ordering within any single topic is strictly sequential and deterministic: each message is assigned an incremental, monotonically increasing offset (`td.Topic.Messages[offset]`).

3. Dedicated Subscriber Worker Routines:
Every active consumer is allocated its own dedicated background worker goroutine (`subscriberWorker`):
```go
go td.subscriberWorker(client, clients, notify, done, stop)
```
This architecture provides critical isolation:
- Slow Subscriber Protection: A lagging or slow consumer reading from its FIFO does not block other consumers or delay message publication on the topic.
- Non-Blocking Notifications: When a new message is appended to a topic, the topic distributor filters matching clients and sends a signal across a 1-slot buffered channel (`client.Notify`). If the subscriber worker is already active or busy delivering a previous batch, the notification does not block.
- Batch Catch-Up: When waking up, the subscriber worker reads all available messages from its last processed offset up to the latest topic message snapshot.

### Synchronization Strategy and Shared State

Shared state is protected through a combination of read-write locks, atomic primitives, and channel coordination:

- Topic State Synchronization (`sync.RWMutex`):
The `Topics` container protects the global topic registry with a `sync.RWMutex`, allowing concurrent reads during topic lookups while serializing new topic creations. Inside each topic (`TopicDatas`), an independent `sync.RWMutex` protects the message slice and subscriber prefix tree.
- Append-Only Snapshotting:
Messages in a topic are append-only and immutable. When a subscriber worker prepares to deliver messages, it acquires a short lock to slice a snapshot (`td.Topic.Messages[next:]`) and immediately releases the lock. Message transmission over IPC occurs entirely outside of locks, preventing slow network or pipe I/O from causing lock contention.
- Atomic Offsets:
Consumer offsets are stored as atomic uint32 variables (`client.Offset.Load()`). Acknowledgments and progress updates are lock-free and thread-safe.
- Graceful Termination:
Shutdown is coordinated through context cancellation (`context.Context`), broadcast channels (`close(td.Done)`), and wait groups (`sync.WaitGroup`). A 5-second timeout (`ShutdownTimeout`) enforces timely cleanup if external processes fail to respond.

## IPC Choice

### Justification of Named Pipes (FIFOs)

The subject allows choosing between System V message queues, POSIX message queues, and Named Pipes (FIFOs). Named Pipes were selected based on the following technical considerations:

1. Stream Abstraction and Elimination of Size Constraints:
System V and POSIX message queues enforce rigid kernel-level quotas on queue depth and individual message size (e.g. `/proc/sys/fs/mqueue/msg_max` and `msgsize_max`, typically limited to 10 messages of 8KB by default on unprivileged Linux kernels). In contrast, Named Pipes operate as standard byte streams backed by kernel pipe buffers. This allows arbitrary message sizes and variable-length streaming without requiring elevated system privileges or hitting artificial kernel limits.

2. Resource Lifecycle and Leak Prevention:
System V message queues are notoriously problematic because they persist in kernel memory even after process crashes, requiring manual cleanup with `ipcs` and `ipcrm`. Named Pipes reside as filesystem nodes (e.g., in `/tmp`). Their lifetime is cleanly managed via standard filesystem operations (`os.Remove`), and they are automatically unlinked upon server shutdown or signal interception.

3. Native Go Integration and Non-Blocking I/O:
In Go, Named Pipes map directly to standard `*os.File`, `bufio.Reader`, and `bufio.Writer` interfaces. This enables idiomatic use of `select` statements, non-blocking open operations (`syscall.O_NONBLOCK`), and I/O write timeouts (`SubscriberWriteTimeout`), which would otherwise require complex CGo wrappers when using POSIX or System V message queue APIs.

### Bidirectional Communication

Unix FIFOs are fundamentally unidirectional byte channels (one writing end, one reading end). Full bidirectional communication between the server and multiple concurrent clients is achieved through the following dual-channel design:

1. Inbound Channel (Client to Server):
The server creates a well-known public FIFO at `/tmp/operator.server.<Server_PID>`. All clients open this FIFO in write-only mode (`os.O_WRONLY | syscall.O_NONBLOCK`) to transmit commands and published messages to the server.

2. Outbound Channel (Server to Client):
When a client starts, it creates its own private response FIFO named:
```text
<server_entrypoint>.<Client_PID>
```
(for example: `/tmp/operator.server.1234.5678`).
The client opens this private FIFO for reading (`os.O_RDWR` to prevent blocking on open).

3. Handshake and Message Delivery:
When a client sends a command to the server, it transmits its PID and configuration in the payload. The server opens the client's dedicated response FIFO for writing.
- For synchronous commands (`create`, `list`, `info`), the server writes the response JSON back to the client's FIFO and closes the writer.
- For asynchronous subscribers (`subscribe`), the server maintains an active writer handle to the client's dedicated FIFO, streaming matched messages continuously as they are published.

This separation guarantees that client responses and subscriber streams never collide or interleave on a single channel.

## Data Structures

### Custom Client Metadata Hashmap

The subject requires a custom-built hashmap with collision resolution for storing client metadata; standard language built-in maps are strictly prohibited for client indexing.

- Collision Resolution: Separate Chaining. The hashmap is structured as an array of 64 bucket pointers (`[]*Node`). Each bucket points to a singly-linked list of `Node` structs storing the key, client pointer, and next pointer.
- Hash Function: Fowler-Noll-Vo 64-bit non-cryptographic hash algorithm (FNV-1a):
  - Offset Basis: `14695981039346656037` (`FnvOffsetBasis`)
  - Prime: `1099511628211` (`FnvPrime`)
  The algorithm iterates through the bytes of the client identifier, XORing the hash with each byte and multiplying by the FNV prime, followed by modulo reduction over the bucket count. FNV-1a provides excellent dispersion, minimal clustering, and fast execution for alphanumeric strings.
- Concurrency and Safety: The hashmap is protected by an internal `sync.RWMutex`. Lookups (`Get`, `Contains`, `Len`) acquire a read lock (`RLock`), while modifications (`Append`, `Delete`) acquire an exclusive lock (`Lock`).
- Complexity: Average-case time complexity is O(1) for lookups, insertions, and deletions; worst-case complexity is O(k) where k is the length of a single collision chain.

### Prefix-Matching Structure: Trie (Prefix Tree)

To deliver messages to subscribers based on key prefixes, a character-based Trie (Prefix Tree) was implemented.

- Structure Design:
  Each node in the Trie represents a character in a prefix:
  ```go
  type Tree struct {
      Char       byte
      Children   map[byte]*Tree
      Clients    []*pr.Client
      AllClients []*pr.Client
      mu         sync.RWMutex
  }
  ```
  - `Children`: Map of character transitions to child Trie nodes.
  - `Clients`: Slice of subscribers whose filter prefix terminates at this exact node.
  - Root Node (`Char = 0`): Stores wildcard subscribers registered with an empty prefix (`""`).

- Justification against Alternative Data Structures:
  - Linear Array / Slice: Checking every consumer's prefix against a message key requires O(N * M) string comparisons per published message, where N is the number of subscribers and M is average prefix length. This degrades severely under high subscriber counts.
  - Binary Search Trees (AVL, Red-Black): Standard balanced search trees sort full strings lexicographically. Performing a prefix range query requires O(log N + K) comparisons, but finding all matching prefix prefixes for a specific message key requires multiple lookups or complex range traversals.
  - Trie (Chosen):
    - Lookup time depends strictly on message key length K, giving O(K) lookup complexity regardless of how many thousands of subscribers are registered.
    - Single-pass filtering: As the server descends the Trie following the characters of the message key, it immediately collects all subscribers attached to each intermediate prefix node.
    - Direct wildcard handling: Subscribers with an empty prefix (`""`) are stored at the root and retrieved in O(1) without tree traversal.

## Testing

The prefix-matching functionality is verified through a comprehensive unit test suite written with Go's native `testing` framework.

### Test Suites Overview

The prefix unit tests are located in `src/tests/prefix/` and cover three main categories:

1. Basic Functionality (`basic_test.go`):
- Exact matches (e.g. prefix `user` matches message key `user`).
- Hierarchical prefix matches (e.g. prefix `user.` matches `user.signup` and `user.profile.update`).
- Non-matching keys (e.g. prefix `user.` does not match `admin.login`).
- Dynamic subscriber management: insertion, updating a subscriber's prefix, and subscriber removal.

2. Edge Cases (`edge_case_test.go`):
- Wildcard subscriptions: empty prefix (`""`) receiving all messages.
- Single-character keys and prefixes (`a`, `b`).
- Duplicate client insertions and removal of non-existent clients.
- Deep prefix hierarchies and special characters (periods, hyphens, underscores, accents, UTF-8 strings).

3. Scalability and Stress Testing (`big_scale_test.go`):
- Scale: 10,000 synthetic subscribers registered with random prefixes generated from realistic word lists and separators.
- Volume: 3,000 distinct message keys evaluated against the Trie.
- Validation: Every lookup is cross-checked against a linear reference model to guarantee 100% filtering accuracy.
- Thread Safety: Concurrent insertions, removals, and lookups executed simultaneously across multiple goroutines to verify race freedom.

### Running the Tests

- Run all project test suites (parsing, hashmap, prefix, ipc, pipeline):
```bash
make test
```

- Run all tests with the Go race detector enabled:
```bash
make test-race
```

- Run given binary tests:
```bash
make test-bubble
```

- Run only the prefix-matching unit tests:
```bash
go test ./src/tests/prefix/... -v
```

- Run prefix unit tests with race detection and test coverage:
```bash
go test -race -v -cover ./src/tests/prefix/...
```

- Run full project verification (git checks, formatting, linting, deadcode, full build, strict race tests):
```bash
make check
```

## Resources

- [IPC system video](https://www.youtube.com/watch?v=Y2mDwW2pMv4)

- In order to fully understand the main concepts and tools to use, we used NotebookLM like a teacher, giving him multiple sources:
	- [Project subject](https://cdn.intra.42.fr/pdf/pdf/229762/en.subject.pdf) to understand each needed part
	- [Interprocess Communication PDF](https://pdfcoffee.com/ipc-basic-ppt-pdf-free.html) to understand basic FIFO system
	- [Project git guideline](https://cdn.intra.42.fr/pdf/pdf/229761/en.git_guidelines.pdf) and [git/github tuto](https://pdfcoffee.com/complete-git-and-github-tutorial-for-beginners-ys-pdf-free.html) to understand project management

- To help us with terminology and project management, we created little pdf thanks to NotebookLM:
	- [Git management](git-branch-management.pdf)
	- [Naming conventions](git-naming-conventions.pdf)
	- [GO Tests](go-tests-guide.pdf)
	- [IPC Protocol](ipc-protocol-guide.pdf)
	- [Project Terminology](treenity-terminology.pdf)
