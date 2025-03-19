package helper

import (
	"log"
	"os"
	"path/filepath"
)

func GetAllFiles(FolderPath string) []string {
	files := make([]string, 0)
	err := filepath.Walk(FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return files
}
