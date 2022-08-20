package mate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	postsaver "github.com/mpragliola/mate/internal/postsaver"
)

type Writer struct {
	postSaver postsaver.PostSaver
}

func NewWriter(ps postsaver.PostSaver) *Writer {
	return &Writer{
		postSaver: ps,
	}
}

func (w *Writer) Write(post Post, project *Project) error {
	destinationPath := filepath.Join(
		project.GetPublicDirectory(),
		post.Path,
	)
	if _, err := os.Stat(destinationPath); err != nil {
		os.MkdirAll(destinationPath, 0755)
	}

	fileName := post.File + ".html"

	layoutPath := layoutPathProvider(post.Layout, project)
	templatedSource, err := parseLayout(layoutPath, post, project)
	if err != nil {
		return err
	}

	err = w.postSaver.Save(filepath.Join(destinationPath, fileName), templatedSource)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WritePages(posts []Post, project *Project) error {
	wg := sync.WaitGroup{}

	for _, p := range posts {
		wg.Add(1)
		go func(p Post) {
			defer wg.Done()

			w.Write(p, project)
		}(p)
	}

	wg.Wait()

	return nil
}

func (w *Writer) Clean(project *Project) error {
	return os.RemoveAll(project.GetPublicDirectory())
}

func layoutPathProvider(layout string, project *Project) string {
	layoutPath := filepath.Join(
		project.GetLayoutsDirectory(),
		layout+".tpl.html",
	)

	return layoutPath
}

func parseLayout(layoutPath string, post Post, project *Project) (string, error) {
	t, err := template.New(filepath.Base(layoutPath)).Funcs(
		template.FuncMap{
			"tags": func() []string {
				return project.GetTags()
			},
			"linktag": func(tag string) string {
				return filepath.Join(
					project.GetPublicTagsPath(),
					tag+".html",
				)
			},
			"pagessorted": func(mode string) []Post {
				switch mode {
				case "created":
					return project.PostsSorted(project.CreatedOnSorter())
				default:
					return project.Posts()
				}
			},
		},
	).ParseFiles(layoutPath)
	if err != nil {
		return "", fmt.Errorf("layout %s not found", layoutPath)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, post)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
