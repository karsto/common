package files

import (
	"fmt"
	"os"
	"path"
)

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func WriteFile(dest, name, content string) error {
	err := CreateIfNotExists(dest, 0777) // TODO:
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(dest, name))
	if err != nil {
		return err
	}
	_, err = f.WriteString(content)
	if err != nil {
		err := f.Close()
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
