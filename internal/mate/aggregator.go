package mate

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Common markdown extension
const markdownExt = ".md"

type Aggregator struct {
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// AggregatePostPaths scans the posts folder recursively, gathering the markdown files,
// providing an array of their paths
func (a *Aggregator) AggregatePostPaths(p *Project) ([]string, error) {
	paths := make([]string, 0)

	if _, err := os.Stat(p.GetPostsDirectory()); err != nil {
		return nil, fmt.Errorf("posts path [%s] not found", p.GetPostsDirectory())
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
