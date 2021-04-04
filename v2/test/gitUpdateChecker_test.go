package test

import (
	"github.com/DarkFighterLuke/gitUpdateChecker/v2"
	"log"
	"testing"
	"time"
)

// Tests setting repository information
func TestSetRepoInfo(t *testing.T) {
	success := 0
	err := gitUpdateChecker.SetRepoInfo("https://github.com/DarkFighterLuke/test.git", "")
	if err != nil {
		success++
	}

	err = gitUpdateChecker.SetRepoInfo("", "master")
	if err != nil {
		success++
	}

	err = gitUpdateChecker.SetRepoInfo("https://github.com/DarkFighterLuke/test.git", "master")
	if err != nil {
		t.Errorf("An error occurred: %s", err)
	} else {
		success++
	}

	if success != 3 {
		log.Printf("%d/3\n", success)
		t.Fail()
	}
}

// Tests the monitoring process launcher
func TestStartUpdateProcess(t *testing.T) {
	err := gitUpdateChecker.SetRepoInfo("https://github.com/DarkFighterLuke/test.git", "master")
	if err != nil {
		t.Errorf("An error occurred: %s", err)
		return
	}
	ch, stop := gitUpdateChecker.StartUpdateProcess(5 * time.Second)

	go updateMonitor(ch)

	time.Sleep(1 * time.Minute)
	stop <- true
}

func updateMonitor(ch chan bool) {
	for range ch {
		log.Println("New commit found in the given branch!")
	}
}
