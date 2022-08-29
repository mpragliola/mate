package mate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	postsaver "github.com/mpragliola/mate/internal/postsaver"
	"golang.org/x/sync/errgroup"
)

type Writer struct {
	postSaver postsaver.PostSaver
}

func NewWriter(ps postsaver.PostSaver) *Writer {
	return &Writer{
		postSaver: ps,
	}
}

func (w *Writer) WriteTag(tag string, project *Project) error {
	destinationPath := filepath.Join(
		project.GetPublicDirectory(),
		project.GetPublicTagsPath(),
	)
	if _, err := os.Stat(destinationPath); err != nil {
		os.MkdirAll(destinationPath, 0755)
	}

	fileName := tag + ".html"

	layoutPath := layoutPathProvider("tag", project)

	templatedSource, err := parseLayout(layoutPath, tag, project)
	if err != nil {
		return err
	}

	err = w.postSaver.Save(filepath.Join(destinationPath, fileName), templatedSource)
	if err != nil {
		return err
	}

	return nil
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
	templatedSource, err := parseLayout(layoutPath, &post, project)
	if err != nil {
		return err
	}

	err = w.postSaver.Save(filepath.Join(destinationPath, fileName), templatedSource)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteTags(posts []Post, project *Project) error {
	g := new(errgroup.Group)

	for _, t := range project.GetTags() {
		t := t
		g.Go(func() error {
			return w.WriteTag(t, project)
		})
	

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WritePages(posts []Post, project *Project) error {
	g := new(errgroup.Group)

	for _, p := range posts {
		p := p
		g.Go(func() error {
			return w.Write(p, project)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

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

func parseLayout(layoutPath string, data interface{}, project *Project) (string, error) {
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
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
