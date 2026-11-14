package ipc

import (
	"fmt"
	"testing"
	pr "treenity/src/protocol"
	ipc "treenity/src/protocol/ipc"
	sce "treenity/src/tests/scenarios"
)

type IpcWorker struct {
	Writer *ipc.Writer
	Reader *ipc.Reader
}

func UtilTestAction(t *testing.T, ipcMap map[int]*IpcWorker, step sce.IpcStep) {
	switch step.Action {
	case sce.IpcActionCreateReader:
		// Create Reader
		new_reader, reader_err := ipc.NewReader(get_path(step.Pid))
		if (reader_err != nil) != step.WantErr {
			t.Errorf("NewReader() error = %v, WantErr %v", reader_err, step.WantErr)
		}

		// Stock Reader
		cur_worker := ipcMap[step.Pid]
		if cur_worker == nil {
			cur_worker = &IpcWorker{
				Reader: new_reader,
			}
			ipcMap[step.Pid] = cur_worker
		} else {
			cur_worker.Reader = new_reader
		}

	case sce.IpcActionCreateWriter:
		// Create Writer
		new_writer, writer_err := ipc.NewWriter(get_path(step.Pid))
		if (writer_err != nil) != step.WantErr {
			t.Errorf("NewWriter() error = %v, WantErr %v", writer_err, step.WantErr)
		}

		// Stock Writer
		cur_worker := ipcMap[step.Pid]
		if cur_worker == nil {
			cur_worker = &IpcWorker{
				Writer: new_writer,
			}
			ipcMap[step.Pid] = cur_worker
		} else {
			cur_worker.Writer = new_writer
		}

	case sce.IpcActionRead:
		// Get Reader
		reader := ipcMap[step.Pid].Reader
		if reader == nil && step.Pid != -1 {
			t.Errorf("Reader is a nil pointer")
		}

		// Get text
		text, read_err := reader.ReadLine()

		// Compare result
		if (read_err != nil) != step.WantErr {
			t.Errorf("ReadLine() error = %v, WantErr %v", read_err, step.WantErr)
		} else if string(text) != step.Text && !step.WantErr {
			t.Errorf("ReadLine() got = %s, want %s", string(text), step.Text)
		}

	case sce.IpcActionWrite:
		// Get Writer
		writer := ipcMap[step.Pid].Writer
		if writer == nil && step.Pid != -1 {
			t.Errorf("Writer is a nil pointer")
		}

		// Write text
		write_err := writer.WriteLine([]byte(step.Text))

		// Compare result
		if (write_err != nil) != step.WantErr {
			t.Errorf("WriteLine() error = %v, WantErr %v", write_err, step.WantErr)
		}
	}
}

func get_path(pid int) string {
	return fmt.Sprintf("%s%d", pr.EntryPoint, pid)
}
