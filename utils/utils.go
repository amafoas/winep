package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func AbsFolderPath(path string) (string, error) {
	folderInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}
	if !folderInfo.IsDir() {
		return "", fmt.Errorf("the path is invalid or does not correspond to an existing folder")
	}
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	return pathAbs, nil
}

func GetArchFromPrefix(folderPath string) (string, error) {
	path := folderPath + "/system.reg"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile("#arch=.+")
	archString := r.FindString(string(fileBytes))
	arch := strings.Split(archString, "=")

	var res string
	if len(arch) < 2 {
		res = "unknown arch"
	}
	if arch[1] == "win32" {
		res = "32 bits"
	} else if arch[1] == "win64" {
		res = "64 bits"
	}

	return res, nil
}

func isValidID(id string) error {
	if len(id) < 2 {
		return fmt.Errorf("id must be at least 2 characters")
	}
	//onli valid [aA-zZ] -
	if id == "az" {
		return fmt.Errorf("id only can have alphanumeric character and - _")
	}

	if id == "-_" {
		return fmt.Errorf("id must start with an alphanumeric character")
	}

	if id == "-_" {
		return fmt.Errorf("id cannot end with special characters")
	}

	return nil
}
