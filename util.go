package suda

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func intn(v int) int {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rd.Intn(v)
}

func loadYaml(name string, data interface{}) error {
	b, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(b, data); err != nil {
		return err
	}

	return nil
}

func readDirFileInfo(name string) ([]fs.FileInfo, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func readDirNames(name string) ([]string, error) {
	names := []string{}
	infos, err := readDirFileInfo(name)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		names = append(names, info.Name())
	}
	return names, nil
}

func genRequestId() string {
	b := make([]byte, 16)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, err := rd.Read(b)
	if err != nil {
		return ""
	}
	return "req_" + base64.RawURLEncoding.EncodeToString(b)
}

func writeBody(w io.Writer, code int, body string) error {
	r := &http.Response{}
	r.Header = http.Header{}
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
	r.Header.Set("X-Content-Type-Options", "nosniff")
	r.StatusCode = code
	r.Body = io.NopCloser(bytes.NewBufferString(body))
	return r.Write(w)
}
