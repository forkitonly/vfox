package sdk

import (
	"errors"
	util2 "github.com/aooohan/version-fox/util"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Operation struct {
	// eg: ~/.version-fox/.cache/node
	localPath string
	// eg: ~/.version-fox
	vfConfigPath string
	osType       util2.OSType
	archType     util2.ArchType
}

func (s *Operation) Download(url *url.URL) (string, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("source file not found")
	}

	err = os.MkdirAll(s.localPath, 0755)
	if err != nil {
		return "", err
	}

	path := filepath.Join(s.localPath, filepath.Base(url.Path))

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}
