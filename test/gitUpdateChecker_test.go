package test

import (
	"gitUpdateChecker"
	"testing"
	"time"
)

func TestStartUpdateProcess(t *testing.T) {
	gitUpdateChecker.SetRepoInfo("https://github.com/DarkFighterLuke/test.git")
	_, err := gitUpdateChecker.StartUpdateProcess(1 * time.Second)
	if err != nil {
		t.Errorf("An error occurred: %s", err)
	}
}
