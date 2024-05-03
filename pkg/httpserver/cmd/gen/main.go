package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed data/src
var srcFiles embed.FS

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

		if err := render(path.Join("src", dirname, *filename+".go"), tpl, templateVal); err != nil {
			panic(err)
		}
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
