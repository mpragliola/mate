package mate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
)

type Parser struct {
	mu sync.Mutex
}

func NewParser() *Parser {
	return &Parser{}
}

// ParsePaths will scan a slice of local, relative paths (provided ideally by mate.Aggregator)
// performing ParsePath on each one and returning an array of pages.
func (p *Parser) ParsePaths(paths []string, project *Project) ([]Post, error) {
	posts := make([]Post, 0)

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			if post, err := p.ParsePath(path, project); err == nil {
				p.mu.Lock()
				posts = append(posts, *post)
				p.mu.Unlock()
			}

			wg.Done()
		}(path)
	}

	wg.Wait()

	posts = project.PostsSorted(project.CreatedOnSorter())

	return posts, nil
}

func (p *Parser) ParsePath(path string, project *Project) (*Post, error) {
	relativePath, _ := filepath.Rel(project.GetPostsDirectory(), path)

	source, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file " + path)
	}

	fm := NewFrontMatterFromSource(string(source))

	var buf bytes.Buffer
	err = goldmark.Convert([]byte(fm.GetBody()), &buf)
	if err != nil {
		return nil, fmt.Errorf("error parsing file " + path)
	}

	fileName := strings.TrimSuffix(filepath.Base(relativePath), filepath.Ext(relativePath))
	tags := fm.GetArrayValue("tags")

	p.mu.Lock()
	project.AddTags(tags...)
	p.mu.Unlock()

	post := &Post{
		CreatedOn: time.Now(),
		Title:     fm.GetValue("title", fileName),
		Tags:      tags,
		Path:      filepath.Dir(relativePath),
		File:      fileName,
		Source:    fm.GetBody(),
		Body:      template.HTML(buf.String()),
		Layout:    fm.GetValue("layout", "post"),
		Slug:      fileName,
		Project:   project,
	}

	p.mu.Lock()
	project.AddPost(*post)
	p.mu.Unlock()

	return post, nil
}
