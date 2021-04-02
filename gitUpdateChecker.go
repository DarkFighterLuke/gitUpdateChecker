package gitUpdateChecker

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"log"
	"time"
)

type RepoInfo struct {
	address string
}

var info RepoInfo

func SetRepoInfo(repoAddress string) {
	info.address = repoAddress
}

func StartUpdateProcess(timeInterval time.Duration) (chan bool, error) {
	if info.address == "" {
		return nil, fmt.Errorf("empty repository address")
	}
	ch := make(chan bool)
	go updateChecker(ch, timeInterval)

	return ch, nil
}

func updateChecker(ch chan bool, timeInterval time.Duration) {
	var oldHash string
	var isFirstRun = true

	for {
		r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: info.address,
		})
		if err != nil {
			log.Printf("An error occurred while cloning repository: %s", err)
			return
		}
		ref, err := r.Head()
		if err != nil {
			log.Printf("Could not get current HEAD: %s\nRetrying...", err)
			continue
		}
		hash := ref.Hash().String()
		if isFirstRun {
			isFirstRun = false
			oldHash = hash
		} else if hash != oldHash {
			ch <- true
		}
		time.Sleep(timeInterval)
	}
}
