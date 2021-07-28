package osutil

import (
	"os"
)

func Exists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsSymlink(filePath string) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}
