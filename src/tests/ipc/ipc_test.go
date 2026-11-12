package ipc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"treenity/src/protocol/ipc"
	sce "treenity/src/tests/scenarios/ipc"
)

func TestNewWriterWithoutFifo(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.fifo")
	var newWriter *ipc.Writer
	var err error

	newWriter, err = ipc.NewWriter(path)
	if err == nil {
		newWriter.Close()
		t.Fatalf("%v", err)
	}
}

func TestFifoLifeCycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.fifo")

	newReader, err := ipc.NewReader(path)
	if err != nil {
		t.Fatalf("%v", err)
	}

	info, err := os.Stat(path)
	if err != nil && info.Mode()&os.ModeNamedPipe != 0 {
		t.Fatalf("%v", err)
	}
	
	newWriter, err := ipc.NewWriter(path)
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = newWriter.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}

	err = newReader.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}
	
	_, err = os.Stat(path) 
	if !os.IsNotExist(err) {
		t.Fatalf("fifo still exist: %v", err)
	}
}

func TestIpc(t *testing.T) {
	for _, tt := range sce.IpcTestScenario {
		t.Run(tt.Name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.fifo")

			newReader, err := ipc.NewReader(path)
			if err != nil {
				t.Fatalf("%v", err)
			}

			newWriter, err := ipc.NewWriter(path)
			if err != nil {
				t.Fatalf("%v", err)
			}

			err = newWriter.WriteLine([]byte(tt.Input))
			readPayload, err := newReader.ReadLine()
			
			if !bytes.Equal(readPayload, []byte(tt.Want))  {
				t.Errorf("Payload = %v, want %v", readPayload, tt.Want)
			}
		})
	}
}




