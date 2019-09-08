package util

import (
	"glt/com.weiai.www/common/log"
	"os"
	"path/filepath"
)

func ReadFile(path string) *os.File {

	f, err := os.Open(path)
	if err != nil {
		panic("open failed!")
	}
	return f
}

func getFilelist(path string) (files []string) {

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Error("filepath.Walk() returned %v\n", err)
	}
	return files
}
