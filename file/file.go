package file

import (
	"fmt"
	"os"
)

func CreateDir(path string) error {
	if PathExists(path) {
		return nil
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// create file with dir if dir is not exist
// path is dir
// name is file name
func CreateFileWithDir(path string, name string, content string) error {
	if err := CreateDir(path); err != nil {
		return err
	}
	file, _ := os.OpenFile(fmt.Sprintf("%v/%v", path, name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}
