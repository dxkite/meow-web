package main

import (
	"archive/zip"
	"bytes"
	"embed"
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

//go:embed data/tpl
var srcFiles embed.FS

//go:embed data/pkg.zip
var pkgFile []byte

//go:embed data/go.mod.tpl
var goModStr string

//go:embed data/go.sum
var goSumStr string

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

	if err := renderString(goModStr, templateVal, false, path.Join(*output, "go.mod")); err != nil {
		panic(err)
	}

	if err := renderString(goSumStr, templateVal, false, path.Join(*output, "go.sum")); err != nil {
		panic(err)
	}

	if err := fs.WalkDir(srcFiles, "data/tpl/entity", func(p string, d fs.DirEntry, err error) error {
		fmt.Println("scan file", p)

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go.tpl") {
			return nil
		}

		dirname := strings.TrimSuffix(d.Name(), ".go.tpl")

		tplStr, err := srcFiles.ReadFile(p)
		if err != nil {
			return err
		}

		outFile := path.Join(*output, "src", dirname, *filename+".go")

		fmt.Println("prepare file", p, d.Name(), "-->", outFile)

		if err := renderString(string(tplStr), templateVal, false, outFile); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := fs.WalkDir(srcFiles, "data/tpl/pkg", func(p string, d fs.DirEntry, err error) error {
		fmt.Println("scan file", p)

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go.tpl") {
			return nil
		}

		outputFilename := strings.TrimSuffix(strings.TrimPrefix(p, "data/tpl"), ".tpl")

		tplStr, err := srcFiles.ReadFile(p)
		if err != nil {
			return err
		}

		outFile := path.Join(*output, outputFilename)

		fmt.Println("prepare file", p, d.Name(), "-->", outFile)

		if err := renderString(string(tplStr), templateVal, false, outFile); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := unzip(pkgFile, path.Join(*output, "pkg")); err != nil {
		panic(err)
	}
}

func renderString(tplStr string, val interface{}, overwrite bool, p string) error {
	tpl, err := template.New("template").Parse(tplStr)
	if err != nil {
		panic(err)
	}

	if !overwrite && exists(p) {
		fmt.Printf("file %s is exist, deleted to generate\n", p)
		return nil
	}

	dir := filepath.Dir(p)
	os.MkdirAll(dir, os.ModePerm)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	fmt.Printf("file %s is created\n", p)
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
