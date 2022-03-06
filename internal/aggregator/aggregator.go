package aggregator

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/mpragliola/mate/internal/mate"
)

// Common markdown extension
const (
	markdownExt              = ".md"
	errorMessagePathNotFound = "posts path [%s] not found"
)

// AggregatePostPaths scans the posts folder recursively, gathering the markdown files,
// providing an array of their paths
func AggregatePostPaths(p *mate.Project) ([]string, error) {
	paths := make([]string, 0)

	if _, err := os.Stat(p.GetPostsDirectory()); err != nil {
		return nil, fmt.Errorf(errorMessagePathNotFound, p.GetPostsDirectory())
	}

	err := filepath.WalkDir(
		p.GetPostsDirectory(),
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != markdownExt {
				return nil
			}

			paths = append(paths, path)

			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return paths, nil
}
