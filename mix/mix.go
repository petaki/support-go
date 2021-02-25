package mix

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Mix type.
type Mix struct {
	url       string
	manifests map[string]map[string]string
}

// New function.
func New(url string) *Mix {
	m := new(Mix)
	m.url = url
	m.manifests = make(map[string]map[string]string)

	return m
}

// Mix function.
func (m *Mix) Mix(path, manifestDirectory string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if manifestDirectory != "" && !strings.HasPrefix(manifestDirectory, "/") {
		manifestDirectory = "/" + manifestDirectory
	}

	_, err := os.Stat("./public" + manifestDirectory + "/hot")
	if os.IsExist(err) {
		content, err := ioutil.ReadFile("./public" + manifestDirectory + "/hot")
		if err != nil {
			return "", err
		}

		url := strings.TrimSpace(string(content))

		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return url[strings.Index(url, ":")+1:] + path, nil
		}

		return "//localhost:8080" + path, nil
	}

	manifestPath := "./public" + manifestDirectory + "/mix-manifest.json"

	if _, ok := m.manifests[manifestPath]; !ok {
		_, err := os.Stat(manifestPath)
		if os.IsNotExist(err) {
			return "", errors.New("mix: the mix manifest does not exist")
		}

		content, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return "", err
		}

		var data map[string]string

		err = json.Unmarshal(content, &data)
		if err != nil {
			return "", err
		}

		m.manifests[manifestPath] = data
	}

	manifest := m.manifests[manifestPath]

	if _, ok := manifest[path]; !ok {
		return "", fmt.Errorf("mix: unable to locate mix file: %v", path)
	}

	return m.url + manifestDirectory + manifest[path], nil
}

// Hash function.
func (m *Mix) Hash(manifestDirectory string) (string, error) {
	if manifestDirectory != "" && !strings.HasPrefix(manifestDirectory, "/") {
		manifestDirectory = "/" + manifestDirectory
	}

	manifestPath := "./public" + manifestDirectory + "/mix-manifest.json"

	_, err := os.Stat(manifestPath)
	if os.IsNotExist(err) {
		return "", errors.New("mix: the mix manifest does not exist")
	}

	file, err := os.Open(manifestPath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
