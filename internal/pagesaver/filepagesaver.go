package pagesaver

import (
	"os"
)

type FilePageSaver struct {
}

func (fps *FilePageSaver) Save(file, content string) error {
	err := os.WriteFile(file, []byte(content), 0644)

	return err
}
