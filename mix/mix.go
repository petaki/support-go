package mix

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

// Mix type.
type Mix struct {
	url         string
	publicPath  string
	hotProxyURL string
	manifests   map[string]map[string]string
}

// New function.
func New(url, publicPath, hotProxyURL string) *Mix {
	m := new(Mix)
	m.url = url
	m.publicPath = publicPath
	m.hotProxyURL = hotProxyURL
	m.manifests = make(map[string]map[string]string)

	return m
}

// Mix function.
func (m *Mix) Mix(path, manifestDirectory string) (string, error) {
	path = m.pathPrefix(path)
	manifestDirectory = m.pathPrefix(manifestDirectory)

	_, err := os.Stat(m.publicPath + manifestDirectory + "/hot")
	if os.IsExist(err) {
		if m.hotProxyURL != "" {
			return m.hotProxyURL + path, nil
		}

		content, err := ioutil.ReadFile(m.publicPath + manifestDirectory + "/hot")
		if err != nil {
			return "", err
		}

		url := strings.TrimSpace(string(content))

		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return url[strings.Index(url, ":")+1:] + path, nil
		}

		return "//localhost:8080" + path, nil
	}

	manifestPath := m.publicPath + m.manifestPath(manifestDirectory)

	if _, ok := m.manifests[manifestPath]; !ok {
		_, err := os.Stat(manifestPath)
		if os.IsNotExist(err) {
			return "", ErrManifestNotExist
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
	manifestPath := m.publicPath + m.manifestPath(m.pathPrefix(manifestDirectory))

	_, err := os.Stat(manifestPath)
	if os.IsNotExist(err) {
		return "", ErrManifestNotExist
	}

	file, err := os.Open(manifestPath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return m.hashFromFile(file)
}

// HashFromFS function.
func (m *Mix) HashFromFS(manifestDirectory string, staticFS fs.FS) (string, error) {
	file, err := staticFS.Open(strings.TrimPrefix(m.manifestPath(manifestDirectory), "/"))
	if err != nil {
		return "", err
	}

	defer file.Close()

	return m.hashFromFile(file)
}

func (m *Mix) hashFromFile(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (m *Mix) manifestPath(manifestDirectory string) string {
	return manifestDirectory + "/mix-manifest.json"
}

func (m *Mix) pathPrefix(path string) string {
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}
