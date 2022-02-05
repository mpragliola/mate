package watcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mpragliola/mate/internal/mate"
)

type Watcher struct {
	aggregator mate.Aggregator
	project    mate.Project
	oldPaths   map[string][]byte
}

func NewWatcher(aggregator mate.Aggregator, project mate.Project) Watcher {
	return Watcher{aggregator, project, make(map[string][]byte)}
}

func (w *Watcher) Watch() {
	for {
		paths, _ := w.aggregator.AggregatePostPaths(&w.project)

		changedPaths := make([]string, 0)
		deletedPaths := make([]string, 0)

		for path, _ := range w.oldPaths {
			found := false

			for _, newPath := range paths {
				if path == newPath {
					found = true
					break
				}
			}

			if !found {
				deletedPaths = append(deletedPaths, path)
			}
		}

		for _, path := range paths {
			oldHash, exists := w.oldPaths[path]
			newHash := getChecksum(path)

			if !exists || string(newHash) != string(oldHash) {
				changedPaths = append(changedPaths, path)
				w.oldPaths[path] = newHash
			}
		}

		fmt.Println("Changed paths:", deletedPaths)
		fmt.Println("Deleted paths:", changedPaths)

		time.Sleep(time.Millisecond * 1000)
	}
}

func getChecksum(path string) []byte {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return hash.Sum(nil)
}
