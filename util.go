package suda

import (
	"context"
	"encoding/base64"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dxkite.cn/log"
	"github.com/dlclark/regexp2"
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

func readAuthData(req *http.Request, source []AuthSourceConfig) string {
	for _, v := range source {
		if v.Type == "cookie" {
			if vv, err := req.Cookie(v.Name); err == nil {
				return vv.Value
			}
		}
		if v.Type == "header" {
			if vv := req.Header.Get(v.Name); vv != "" {
				return vv
			}
		}
	}
	return ""
}

func copyHeader(w http.ResponseWriter, h http.Header) {
	for k, v := range h {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
}

func isUpgradeToWebsocket(req *http.Request) bool {
	connection := req.Header.Get("Connection")
	upgrade := req.Header.Get("Upgrade")
	if strings.ToLower(connection) == "upgrade" &&
		strings.ToLower(upgrade) == "websocket" {
		return true
	}
	return false
}

func regexReplaceAll(reg, input, replacement string) (string, error) {
	r, err := regexp2.Compile(reg, 0)
	if err != nil {
		return "", err
	}
	v, err := r.Replace(input, replacement, -1, -1)
	if err != nil {
		return "", err
	}
	return v, nil
}

func applyLogConfig(ctx context.Context, level int, output string) {
	if level != 0 {
		log.SetLevel(log.LogLevel(level))
	}
	if output == "" {
		return
	}
	log.Println("log output file", output)
	filename := output
	var w io.Writer
	if f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err != nil {
		log.Warn("log file open error", filename)
		return
	} else {
		w = f
		if filepath.Ext(filename) == ".json" {
			w = log.NewJsonWriter(w)
		} else {
			w = log.NewTextWriter(w)
		}
		go func() {
			<-ctx.Done()
			_ = f.Close()
		}()
	}
	log.SetOutput(log.MultiWriter(w, log.Writer()))
}
