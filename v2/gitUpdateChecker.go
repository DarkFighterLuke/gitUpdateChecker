package gitUpdateChecker

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
	"time"
)

type RepoInfo struct {
	address string
	branch  string
}

var info RepoInfo

// Sets repo address
func SetRepoInfo(repoAddress, branch string) error {
	if repoAddress == "" {
		return fmt.Errorf("empty repository address")
	}
	if branch == "" {
		return fmt.Errorf("empty branch name")
	}

	info.address = repoAddress
	info.branch = branch

	return nil
}

// Starts the monitoring process and returns a channel to listen for updates
func StartUpdateProcess(timeInterval time.Duration) (chan bool, chan bool) {
	ch := make(chan bool)
	stop := make(chan bool)
	go updateChecker(ch, stop, timeInterval)

	return ch, stop
}

// Monitors repo address checking for new commits and returns true on the channel if there are ones
func updateChecker(ch, stop chan bool, timeInterval time.Duration) {
	var oldHash plumbing.Hash
	var isFirstRun = true

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{info.address},
	})

	for {
		select {
		case s := <-stop:
			if s {
				log.Println("stop")
				return
			}
		default:
			refs, err := remote.List(&git.ListOptions{InsecureSkipTLS: true})
			if err != nil {
				log.Println(err)
				goto sleep
			}

			for _, v := range refs {
				if v.Name().Short() == info.branch {
					if isFirstRun {
						oldHash = v.Hash()
						isFirstRun = false
					} else if v.Hash() != oldHash {
						ch <- true
						oldHash = v.Hash()
					}
				}
			}

		sleep:
			time.Sleep(timeInterval)
		}
	}
}
