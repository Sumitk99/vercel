package helper

import (
	"fmt"
	"github.com/Sumitk99/vercel/upload-service/constants"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"path/filepath"
)

func CloneRepo(GithubUrl, ProjectId string) error {
	currPath, err := os.Getwd()
	directory := filepath.Join(currPath, constants.RepoPath, ProjectId)
	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		URL:      GithubUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	fmt.Println("Repository cloned successfully!")
	return nil
}

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
