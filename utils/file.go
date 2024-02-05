package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func RemoveFileNameExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func GetAllFileInDir(dirPath string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err.Error())
		return nil, err
	}

	return files, nil
}
