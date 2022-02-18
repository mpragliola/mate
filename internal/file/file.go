package file

import (
	"os"
	"path/filepath"

	"github.com/mpragliola/mate/internal/postsaver"
)

func Copy(source, dest string, fs postsaver.PostSaver) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		paths, err := os.ReadDir(source)
		if err != nil {
			return err
		}

		for _, path := range paths {
			err := Copy(
				filepath.Join(source, path.Name()),
				filepath.Join(dest, path.Name()),
				fs,
			)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if _, err := os.Stat(filepath.Dir(dest)); err != nil {
		err := os.MkdirAll(filepath.Dir(dest), 0755)
		if err != nil {
			return err
		}
	}

	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	err = fs.Save(dest, string(content))
	if err != nil {
		return err
	}

	return nil
}
