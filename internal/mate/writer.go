package mate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/mpragliola/mate/internal/pagesaver"
)

type Writer struct {
	pageSaver pagesaver.PageSaver
}

func NewWriter(ps pagesaver.PageSaver) *Writer {
	return &Writer{
		pageSaver: ps,
	}
}

func (w *Writer) Write(page Page, project *Project) error {
	destinationPath := filepath.Join(
		project.GetPublicDirectory(),
		page.Path,
	)
	if _, err := os.Stat(destinationPath); err != nil {
		os.MkdirAll(destinationPath, 0755)
	}

	fileName := page.File + ".html"

	layoutPath := layoutPathProvider(page.Layout, project)
	templatedSource, err := parseLayout(layoutPath, page, project)
	if err != nil {
		return err
	}

	err = w.pageSaver.Save(filepath.Join(destinationPath, fileName), templatedSource)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WritePages(pages []Page, project *Project) error {
	wg := sync.WaitGroup{}

	for _, p := range pages {
		wg.Add(1)
		go func(p Page) {
			defer wg.Done()

			w.Write(p, project)
		}(p)
	}

	wg.Wait()

	return nil
}

func layoutPathProvider(layout string, project *Project) string {
	layoutPath := filepath.Join(
		project.GetLayoutsDirectory(),
		layout+".tpl.html",
	)

	return layoutPath
}

func parseLayout(layoutPath string, page Page, project *Project) (string, error) {
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
		},
	).ParseFiles(layoutPath)
	if err != nil {
		return "", fmt.Errorf("layout %s not found", layoutPath)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, page)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
