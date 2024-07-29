package form

import (
	"mime/multipart"
	"reflect"
	"testing"
)

type Test struct {
	Name  string `form:"name"`
	Bool  bool   `form:"bool"`
	Int   int    `form:"int_val"`
	Value []string
	File  *multipart.FileHeader   `form:"simple_file"`
	Files []*multipart.FileHeader `form:"multi_files"`
}

func TestForm(t *testing.T) {
	input := map[string]Item{}
	input["name"] = StringSlice{"test1", "test2"}
	input["Value"] = StringSlice{"test1", "test2"}
	input["bool"] = StringSlice{"true"}
	input["int_val"] = StringSlice{"10086"}
	file1 := &multipart.FileHeader{Filename: "test1"}
	file2 := &multipart.FileHeader{Filename: "test2"}
	input["simple_file"] = FileSlice{file1, file2}
	input["multi_files"] = FileSlice{file1, file2}

	got := Test{}
	err := MappingForm(input, &got)
	if err != nil {
		t.Error(err)
		return
	}
	want := Test{
		Name:  "test2",
		Value: []string{"test1", "test2"},
		Bool:  true,
		Int:   10086,
		File:  file2,
		Files: []*multipart.FileHeader{file1, file2},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("not equal %v = %v", want, got)
	}
}
