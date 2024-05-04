package main

import (
	"archive/zip"
	"bytes"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed data/src
var srcFiles embed.FS

//go:embed data/pkg.zip
var pkgFile []byte

type TemplateValue struct {
	Pkg         string
	Name        string
	PrivateName string
	URI         string
}

func main() {
	pkg := flag.String("pkg", "", "package name")
	name := flag.String("name", "", "package name")
	filename := flag.String("filename", "", "filename")
	privateName := flag.String("private-name", "", "private name")
	uri := flag.String("uri", "", "uri path")
	output := flag.String("output", ".", "output path")

	flag.Parse()

	if *filename == "" {
		*filename = strings.ToLower(*name)
	}

	if *privateName == "" {
		*privateName = strings.ToLower(*name)
	}

	if *uri == "" {
		*uri = strings.ToLower(*name) + "s"
	}

	fmt.Println("create entity", *pkg, *name, *privateName, *uri)

	templateVal := &TemplateValue{
		Pkg:         *pkg,
		Name:        *name,
		PrivateName: *privateName,
		URI:         *uri,
	}

	srcList, err := srcFiles.ReadDir("data/src")
	if err != nil {
		panic(err)
	}

	for _, v := range srcList {
		fmt.Println(v.Name())
		dirname := strings.TrimSuffix(v.Name(), ".go.tpl")

		tplStr, err := srcFiles.ReadFile(path.Join("data/src", v.Name()))
		if err != nil {
			panic(err)
		}
		tpl, err := template.New(v.Name()).Parse(string(tplStr))
		if err != nil {
			panic(err)
		}

		if err := render(path.Join(*output, "src", dirname, *filename+".go"), tpl, templateVal); err != nil {
			panic(err)
		}
	}

	if err := unzip(pkgFile, path.Join(*output, "pkg")); err != nil {
		panic(err)
	}
}

func render(p string, tpl *template.Template, val *TemplateValue) error {
	if exists(p) {
		fmt.Printf("file %s is exist, deleted to generate new\n", p)
		return nil
	}

	dir := filepath.Dir(p)
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	return tpl.Execute(f, val)
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func unzip(data []byte, output string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	err = fs.WalkDir(zipReader, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if p == "." {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		extract := path.Join(output, p)
		if exists(extract) {
			fmt.Printf("file %s is exist, deleted to generate new\n", extract)
			return nil
		}

		f, err := zipReader.Open(p)
		if err != nil {
			fmt.Printf("open file %s error %v\n", p, err)
			return nil
		}
		defer f.Close()

		dir := filepath.Dir(extract)
		os.MkdirAll(dir, os.ModePerm)

		out, err := os.OpenFile(extract, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("create file %s error %v\n", extract, err)
			return err
		}

		defer out.Close()

		_, err = io.Copy(out, f)
		return err
	})
	return err
}
