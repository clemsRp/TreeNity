package pipeline

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var serverBin, clientBin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "treenity-pipeline-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	serverBin = filepath.Join(dir, "server")
	clientBin = filepath.Join(dir, "client")

	for bin, pkg := range map[string]string{serverBin: "treenity/src/cmd/server", clientBin: "treenity/src/cmd/client"} {
		if out, err := exec.Command("go", "build", "-o", bin, pkg).CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "build %s: %v\n%s", pkg, err, out)
			os.RemoveAll(dir)
			os.Exit(1)
		}
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func startServer(t *testing.T) (*exec.Cmd, string) {
	t.Helper()

	srv := exec.Command(serverBin)
	stdout, err := srv.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Process.Kill() })

	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil {
		t.Fatalf("cannot read server IPC path: %v", err)
	}
	return srv, strings.TrimSpace(line)
}
