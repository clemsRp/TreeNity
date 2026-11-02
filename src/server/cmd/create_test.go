package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	h "treenity/src/server/hashmap"
	top "treenity/src/server/topics"
)

func TestAcceptNewMessages_DistributionAndBrokenPipe(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "treenity_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	topicsList := top.NewTopics()
	clients := h.CreateHashmap(int(pr.ClientsHashmapSize))
	topicData := topicsList.Create("test_topic")
	go AcceptNewMessages(topicData, clients)

	// Create 2 test FIFOs (one active reader, one that will close to trigger broken pipe)
	sub1Pipe := filepath.Join(tmpDir, "sub1.pipe")
	sub2Pipe := filepath.Join(tmpDir, "sub2.pipe")

	r1, err := ipc.NewReader(sub1Pipe)
	if err != nil {
		t.Fatalf("Failed to create reader 1: %v", err)
	}
	defer r1.Close()

	r2, err := ipc.NewReader(sub2Pipe)
	if err != nil {
		t.Fatalf("Failed to create reader 2: %v", err)
	}

	client1 := &pr.Client{
		Id:       "client1",
		Topic:    "test_topic",
		Offset:   0,
		Prefix:   "sports.",
		IpcPath:  sub1Pipe,
		IsActive: true,
	}

	client2 := &pr.Client{
		Id:       "client2",
		Topic:    "test_topic",
		Offset:   0,
		Prefix:   "",
		IpcPath:  sub2Pipe,
		IsActive: true,
	}

	clients.Append("client1", client1)
	clients.Append("client2", client2)
	topicData.Tree.Insert(client1)
	topicData.Tree.Insert(client2)

	// 1. Send first message to initialize writers for both active clients
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  "sports.football",
		Body: "Goal!",
	}

	// Read from both readers
	line1, err1 := r1.ReadLine()
	if err1 != nil || len(line1) == 0 {
		t.Fatalf("Failed to read from r1: %v", err1)
	}

	line2, err2 := r2.ReadLine()
	if err2 != nil || len(line2) == 0 {
		t.Fatalf("Failed to read from r2: %v", err2)
	}

	// Both clients should now have an active Writer and offset 1
	if client1.Writer == nil || client2.Writer == nil {
		t.Fatalf("Expected writers to be initialized")
	}
	if client1.Offset != 1 || client2.Offset != 1 {
		t.Fatalf("Expected offsets to be 1")
	}

	// 2. Close r2 to trigger broken pipe (EPIPE) on the next write
	_ = r2.Close()

	// Send second message
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  "sports.tennis",
		Body: "Ace!",
	}

	// Read on r1
	line1, err1 = r1.ReadLine()
	if err1 != nil || len(line1) == 0 {
		t.Fatalf("Failed to read second message from r1: %v", err1)
	}

	// Wait briefly for client2 to process the broken pipe
	time.Sleep(50 * time.Millisecond)

	if client2.IsActive {
		t.Errorf("Expected client2 IsActive to be false after EPIPE")
	}
	if client2.Writer != nil {
		t.Errorf("Expected client2 Writer to be nil after disconnect")
	}

	// 3. Test prefix filtering on client1 (which has prefix "sports.")
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  "news.politics",
		Body: "Election",
	}

	time.Sleep(50 * time.Millisecond)

	// client1 offset should remain 2 because "news.politics" does not match "sports."
	if client1.Offset != 2 {
		t.Errorf("Expected client1 offset to be 2, got %d", client1.Offset)
	}
}
