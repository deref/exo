package osutil

import (
	"io/fs"
	"os"

	"golang.org/x/sys/unix"
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

func IsWritable(filePath string) (bool, error) {
	err := unix.Access(filePath, unix.W_OK)
	return err == nil, err
}

// TODO: This might not actually help anything.
// See <https://twitter.com/tomfaulhaber/status/1447376910531448832>.
func WriteFileSync(name string, data []byte, perm fs.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err := f.Sync(); err != nil {
		return err
	}
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
