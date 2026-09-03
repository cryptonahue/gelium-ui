package lib

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentStartHookRefreshesGentleRegistry(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(root, "gentle-ai.log")
	fake := "#!/usr/bin/env bash\nprintf '%s\\n' \"$*\" >> \"$GENTLE_LOG\"\n"
	fakePath := filepath.Join(bin, "gentle-ai")
	if err := os.WriteFile(fakePath, []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}

	script := filepath.Join(repositoryRoot(t), "lib", "scripts", "agent-start.sh")
	cmd := exec.Command("bash", script, root)
	cmd.Env = append(os.Environ(), "PATH="+bin+":"+os.Getenv("PATH"), "GENTLE_LOG="+logPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("agent-start.sh failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "skill registry refreshed") {
		t.Fatalf("hook must report registry refresh, output=%s", output)
	}
	log, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(log) != "skill-registry refresh --force\n" {
		t.Fatalf("unexpected gentle-ai invocation: %q", log)
	}
}
