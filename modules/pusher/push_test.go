package pusher

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/state-alchemists/zaruba/modules/command"
	"github.com/state-alchemists/zaruba/modules/config"
	"github.com/state-alchemists/zaruba/modules/file"
	"github.com/state-alchemists/zaruba/modules/initiator"
)

func TestPush(t *testing.T) {
	baseTestPath := config.GetTestDir()
	testPath := filepath.Join(baseTestPath, "testPush")

	// clone
	cloneCommand := fmt.Sprintf("git clone ssh://git@ztest-git:22/git-server/repos/parent-push.git %s", testPath)
	if _, err := command.RunShellScriptAndReturn(baseTestPath, cloneCommand, []string{}); err != nil {
		t.Errorf("[ERROR] Cannot git clone: %s", err)
	}

	// copy and commit
	if err := file.Copy("../../test-resource/testPush", testPath); err != nil {
		t.Errorf("[ERROR] Cannot copy test-case: %s", err)
		return
	}
	if _, err := command.RunShellScriptAndReturn(testPath, "git commit -am 'first commit'", []string{}); err != nil {
		t.Errorf("[ERROR] Cannot git commit: %s", err)
	}

	// load project config
	p, err := config.CreateProjectConfig(testPath)
	if err != nil {
		t.Errorf("[ERROR] Cannot load config: %s", err)
		return
	}

	if err = initiator.Init(testPath, p); err != nil {
		t.Errorf("[ERROR] Cannot zaruba init: %s", err)
	}

	if _, err = command.RunShellScriptAndReturn(testPath, "git add . -A", []string{}); err != nil {
		t.Errorf("[ERROR] Cannot git add: %s", err)
	}

	if err = Push(testPath, p); err != nil {
		t.Errorf("[ERROR] Cannot zaruba push: %s", err)
	}

}
