package backup

import (
	"os"
	"path/filepath"
)

type FileList struct {
	Files []string
}

func NewFileList() *FileList {
	fl := &FileList{}
	return fl
}

func (f *FileList) Add(path string) {
	f.Files = append(f.Files, path)
}

func (f *FileList) Walk(baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			f.Add(path)
		}

		return nil
	})
}

func (f *FileList) Len() int {
	return len(f.Files)
}
