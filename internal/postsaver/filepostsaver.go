package postsaver

import (
	"os"
)

type FilePostSaver struct {
}

func (fps *FilePostSaver) Save(file, content string) error {
	err := os.WriteFile(file, []byte(content), 0644)

	return err
}
