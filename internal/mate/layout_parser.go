package mate

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

func ParseLayout(layoutPath string, data interface{}, project *Project) (string, error) {
	t, err := template.New(filepath.Base(layoutPath)).Funcs(
		template.FuncMap{
			"tags": func() []Tag {
				return project.GetTags()
			},
			"linktag": func(tag Tag) string {
				return filepath.Join(
					project.GetPublicTagsPath(),
					tag.Name+".html",
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
