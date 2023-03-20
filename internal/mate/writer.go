package mate

import (
	"os"
	"path/filepath"

	postsaver "github.com/mpragliola/mate/internal/postsaver"
	"golang.org/x/sync/errgroup"
)

// Writer ...
type Writer struct {
	postSaver postsaver.PostSaver
}

// NewWriter ...
func NewWriter(ps postsaver.PostSaver) *Writer {
	return &Writer{
		postSaver: ps,
	}
}

// WriteTag ...
func (w *Writer) WriteTag(tag Tag, project *Project) error {
	if _, err := os.Stat(project.GetPublicTagsPath()); err != nil {
		os.MkdirAll(project.GetPublicTagsPath(), 0755)
	}

	fileName := tag.Name + ".html"
	layoutPath := layoutPathProvider("tag", project)

	templatedSource, err := ParseLayout(layoutPath, tag, project)
	if err != nil {
		return err
	}

	err = w.postSaver.Save(filepath.Join(project.GetPublicTagsPath(), fileName), templatedSource)
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
	templatedSource, err := ParseLayout(layoutPath, &post, project)
	if err != nil {
		return err
	}

	err = w.postSaver.Save(filepath.Join(destinationPath, fileName), templatedSource)
	if err != nil {
		return err
	}

	return nil
}

// WriteTags ...
func (w *Writer) WriteTags(posts []Post, project *Project) error {
	g := new(errgroup.Group)

	for _, t := range project.GetTags() {
		t := t
		g.Go(func() error {
			return w.WriteTag(t, project)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// WritePages ...
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

// Clean ...
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
