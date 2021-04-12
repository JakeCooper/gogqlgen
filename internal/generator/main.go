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
	ClientFile   *os.File
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
	clientFile, err := os.Create(fmt.Sprintf("./%s/client.go", path))
	if err != nil {
		panic(err)
	}

	files := []*os.File{objectFile, gqlFile, mutationFile, queryFile, clientFile}

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
			ClientFile:   clientFile,
		},
		Types: make(map[string]string),
	}
}
