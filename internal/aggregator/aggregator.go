package aggregator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Common markdown extension
const (
	markdownExt              = ".md"
	errorMessagePathNotFound = "posts path [%s] not found"
)

// AggregatePostPaths scans a posts folder recursively, gathering the markdown files,
// providing an array of their paths
func AggregatePostPaths(postsDirectory string) ([]string, error) {
	paths := make([]string, 0)

	if _, err := os.Stat(postsDirectory); err != nil {
		return nil, fmt.Errorf(errorMessagePathNotFound, postsDirectory)
	}

	err := filepath.WalkDir(
		postsDirectory,
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
		return nil, err
	}

	return paths, nil
}
