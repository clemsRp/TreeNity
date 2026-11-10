package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	pr "treenity/src/protocol"
	"treenity/src/protocol/ipc"
	"treenity/src/server/cmd"
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
	go cmd.AcceptNewMessages(topicData, clients)

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

	// 1. Send first message to both active clients
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  []byte("sports.football"),
		Body: []byte("Goal!"),
	}

	// Read from both readers
	line1, err1 := r1.ReadLine()
	if err1 != nil || len(line1) == 0 {
		t.Fatalf("Failed to read from r1: %v", err1)
	}
	var msg1 pr.CmdMessageDatas
	if err := json.Unmarshal(line1, &msg1); err != nil {
		t.Fatalf("Expected JSON from r1, got '%s', err: %v", string(line1), err)
	}
	if !bytes.Equal(msg1.Key, []byte("sports.football")) || !bytes.Equal(msg1.Body, []byte("Goal!")) || msg1.Offset != 0 {
		t.Fatalf("Expected sports.football:Goal! offset 0, got %+v", msg1)
	}

	line2, err2 := r2.ReadLine()
	if err2 != nil || len(line2) == 0 {
		t.Fatalf("Failed to read from r2: %v", err2)
	}
	var msg2 pr.CmdMessageDatas
	if err := json.Unmarshal(line2, &msg2); err != nil {
		t.Fatalf("Expected JSON from r2, got '%s', err: %v", string(line2), err)
	}
	if !bytes.Equal(msg2.Key, []byte("sports.football")) || !bytes.Equal(msg2.Body, []byte("Goal!")) || msg2.Offset != 0 {
		t.Fatalf("Expected sports.football:Goal! offset 0, got %+v", msg2)
	}

	// Both clients should now have offset 1
	topicData.Mu.RLock()
	off1 := client1.Offset
	off2 := client2.Offset
	topicData.Mu.RUnlock()
	if off1 != 1 || off2 != 1 {
		t.Fatalf("Expected offsets to be 1, got client1=%d, client2=%d", off1, off2)
	}

	// 2. Close r2 to trigger broken pipe (EPIPE/client closed) on the next write
	_ = r2.Close()

	// Send second message
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  []byte("sports.tennis"),
		Body: []byte("Ace!"),
	}

	// Read on r1
	line1, err1 = r1.ReadLine()
	if err1 != nil || len(line1) == 0 {
		t.Fatalf("Failed to read second message from r1: %v", err1)
	}
	var msg3 pr.CmdMessageDatas
	if err := json.Unmarshal(line1, &msg3); err != nil {
		t.Fatalf("Expected JSON on second message, got '%s', err: %v", string(line1), err)
	}
	if !bytes.Equal(msg3.Key, []byte("sports.tennis")) || !bytes.Equal(msg3.Body, []byte("Ace!")) || msg3.Offset != 1 {
		t.Fatalf("Expected sports.tennis:Ace! offset 1, got %+v", msg3)
	}

	// Wait briefly for client2 to process the broken pipe
	time.Sleep(50 * time.Millisecond)

	topicData.Mu.RLock()
	c2Active := client2.IsActive
	topicData.Mu.RUnlock()
	if c2Active {
		t.Errorf("Expected client2 IsActive to be false after EPIPE")
	}

	// 3. Test prefix filtering on client1 (which has prefix "sports.")
	topicData.Topic.NewMsgChan <- pr.Message{
		Key:  []byte("news.politics"),
		Body: []byte("Election"),
	}

	time.Sleep(50 * time.Millisecond)

	topicData.Mu.RLock()
	off1 = client1.Offset
	topicData.Mu.RUnlock()

	// client1 offset should remain 2 because "news.politics" does not match "sports."
	if off1 != 2 {
		t.Errorf("Expected client1 offset to be 2, got %d", off1)
	}
}

func TestHandleCmdSubscribe_Replay(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "treenity_replay_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	topicsList := top.NewTopics()
	clients := h.CreateHashmap(int(pr.ClientsHashmapSize))
	topicData := topicsList.Create("test_replay")

	// Add 3 messages directly into topic
	topicData.Topic.Messages = append(topicData.Topic.Messages,
		&pr.Message{Key: []byte("sports.golf"), Body: []byte("Hole in one"), Offset: 0},
		&pr.Message{Key: []byte("cooking.cake"), Body: []byte("Delicious"), Offset: 1},
		&pr.Message{Key: []byte("sports.ski"), Body: []byte("Snow"), Offset: 2},
	)

	subPipe := filepath.Join(tmpDir, "sub_replay.pipe")
	reader, err := ipc.NewReader(subPipe)
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}
	defer reader.Close()

	prefix := "sports."
	offset := 0
	cmdSub := pr.CmdSubscribeDatas{
		IpcIdentifier:  filepath.Join(tmpDir, "sub_replay"),
		TopicName:      "test_replay",
		SubscriberName: "pipe",
		Prefix:         &prefix,
		Offset:         &offset,
	}

	_, _, prErr := cmd.HandleCmdSubscribe(cmdSub, topicsList, clients, subPipe)
	if prErr != nil {
		t.Fatalf("HandleCmdSubscribe failed: %v", prErr.Message)
	}
	cmd.ReplayAndInsert(cmdSub, topicsList, clients)

	// Should receive first message: sports.golf
	l1, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("Failed to read first replay message: %v", err)
	}
	var m1 pr.CmdMessageDatas
	if err := json.Unmarshal(l1, &m1); err != nil {
		t.Fatalf("Failed to unmarshal m1: %v", err)
	}
	if !bytes.Equal(m1.Key, []byte("sports.golf")) || m1.Offset != 0 {
		t.Fatalf("Expected sports.golf offset 0, got %+v", m1)
	}

	// Should receive second message: sports.ski (cooking.cake was skipped due to prefix)
	l2, err := reader.ReadLine()
	if err != nil {
		t.Fatalf("Failed to read second replay message: %v", err)
	}
	var m2 pr.CmdMessageDatas
	if err := json.Unmarshal(l2, &m2); err != nil {
		t.Fatalf("Failed to unmarshal m2: %v", err)
	}
	if !bytes.Equal(m2.Key, []byte("sports.ski")) || m2.Offset != 2 {
		t.Fatalf("Expected sports.ski offset 2, got %+v", m2)
	}

	// Check client offset has advanced past all 3 messages
	client, err := clients.Get("pipe")
	if err != nil || client == nil {
		t.Fatalf("Client 'pipe' not found in clients map")
	}
	if client.Offset != 3 {
		t.Fatalf("Expected client offset 3, got %d", client.Offset)
	}
}
