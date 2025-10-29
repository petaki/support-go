package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Hash function.
func Hash(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return md5Hash(file)
}

// HashFromFS function.
func HashFromFS(filePath string, fileFS fs.FS) (string, error) {
	file, err := fileFS.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return md5Hash(file)
}

func md5Hash(file io.Reader) (string, error) {
	hash := md5.New()

	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
