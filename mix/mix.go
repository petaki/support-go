package mix

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/petaki/support-go/file"
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
	if !errors.Is(err, fs.ErrNotExist) {
		if m.hotProxyURL != "" {
			return m.hotProxyURL + path, nil
		}

		content, err := os.ReadFile(m.publicPath + manifestDirectory + "/hot")
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
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrManifestNotExist
		}

		content, err := os.ReadFile(manifestPath)
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

	hash, err := file.Hash(manifestPath)
	if err != nil {
		return "", ErrManifestNotExist
	}

	return hash, nil
}

// HashFromFS function.
func (m *Mix) HashFromFS(manifestDirectory string, assetFS fs.FS) (string, error) {
	manifestPath := strings.TrimPrefix(m.manifestPath(manifestDirectory), "/")

	hash, err := file.HashFromFS(manifestPath, assetFS)
	if err != nil {
		return "", ErrManifestNotExist
	}

	return hash, nil
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
