package file

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Hash function.
func Hash(filePath string) (string, error) {
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

// HashFromContent function.
func HashFromContent(content []byte) string {
	hash, _ := md5Hash(bytes.NewReader(content))

	return hash
}

func md5Hash(reader io.Reader) (string, error) {
	hash := md5.New()

	_, err := io.Copy(hash, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
