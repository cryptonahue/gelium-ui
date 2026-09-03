package lib

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallAgentsInstallsCanonicalRouting(t *testing.T) {
	target := t.TempDir()
	script := filepath.Join(repositoryRoot(t), "lib", "install-agents.sh")
	cmd := exec.Command("bash", script, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("install-agents.sh failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "canonical routing") {
		t.Fatalf("installer must report canonical routing installation, output=%s", output)
	}
	for _, path := range []string{
		filepath.Join(target, "AGENTS.md"),
		filepath.Join(target, "SKILL.md"),
		filepath.Join(target, "skills", "00-agent-routing.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("installer did not install %s: %v", path, err)
		}
	}
}