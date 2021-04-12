package generator

import (
	"fmt"
	"os"

	"github.com/jakecooper/gogqlgen/internal/mapper"
)

type Outs struct {
	GQLFile      *os.File
	ObjectFile   *os.File
	QueryFile    *os.File
	MutationFile *os.File
}

type Generator struct {
	Outs  *Outs
	Types map[string]string
}

var tm = mapper.ToMap

func New(path string, filename string) *Generator {
	os.Mkdir(path, os.FileMode(0755))
	objectFile, err := os.Create(fmt.Sprintf("./%s/types.go", path))
	if err != nil {
		panic(err)
	}
	gqlFile, err := os.Create(fmt.Sprintf("./%s/gql.go", path))
	if err != nil {
		panic(err)
	}
	queryFile, err := os.Create(fmt.Sprintf("./%s/query.go", path))
	if err != nil {
		panic(err)
	}
	mutationFile, err := os.Create(fmt.Sprintf("./%s/mutation.go", path))
	if err != nil {
		panic(err)
	}

	files := []*os.File{objectFile, gqlFile, mutationFile, queryFile}

	for _, file := range files {
		file.Write([]byte("// GENERATED FILE DO NOT EDIT!!!\n\npackage gen\n\n"))
	}
	// Write package name + generated header
	return &Generator{
		Outs: &Outs{
			GQLFile:      gqlFile,
			ObjectFile:   objectFile,
			QueryFile:    queryFile,
			MutationFile: mutationFile,
		},
		Types: make(map[string]string),
	}
}
