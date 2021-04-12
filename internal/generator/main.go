package generator

import (
	"fmt"
	"os"

	"github.com/jakecooper/gogqlgen/internal/mapper"
)

type Generator struct {
	Outfile *os.File
	Types   map[string]string
}

var tm = mapper.ToMap

func New(path string, filename string) *Generator {
	os.Mkdir(path, os.FileMode(0755))
	fp, err := os.Create(fmt.Sprintf("./%s/%s.go", path, filename))
	if err != nil {
		panic(err)
	}
	// Write package name + generated header
	fp.Write([]byte("// GENERATED FILE DO NOT EDIT!!!\n\npackage gen\n\n"))
	return &Generator{
		Outfile: fp,
		Types:   make(map[string]string),
	}
}
