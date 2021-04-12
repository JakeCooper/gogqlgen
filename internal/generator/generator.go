package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/mapper"
	"github.com/jakecooper/gogqlgen/internal/tokens"
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

func (g *Generator) HandleMutation(req interface{}) {
	// UNIMPLEMENTED
}

func (g *Generator) HandleQuery(req interface{}) {
	// UNIMPLEMENTED
}

func (g *Generator) HandleScalar(req interface{}) {
	// UNIMPLEMENTED
}

func (g *Generator) HandleInputObject(req interface{}) {
	// UNIMPLEMENTED
}

func (g *Generator) HandleObject(req interface{}) {
	z := req.(map[string]interface{})
	fmt.Print(z["kind"], " : ", z["name"], "\n")
	fields := z["fields"]
	// Handle Everything Else
	if fields != nil {
		fmt.Print("Fields: ")
		for _, rf := range fields.([]interface{}) {
			field := tm(rf)
			fieldType := tm(field["type"])["ofType"]
			if fieldType != nil {
				fieldTypeName := tm(fieldType)["name"]
				fieldKind := tm(fieldType)["kind"]
				isList := false
				if fieldKind != nil {
					// InternalType
					// TODO Probably recursion IDK
					ofType := tm(fieldType)["ofType"]
					isList = tm(fieldType)["kind"] == "LIST"
					if ofType != nil {
						ofType = tm(ofType)["ofType"]
						if ofType != nil {
							fieldTypeName = tm(ofType)["name"]
						}
					}
				}
				typeName := fieldTypeName.(string)

				if graphql.IsInternal(typeName) {
					// Skip internal GraphQL shit
					continue
				}
				fmt.Print(field["name"], ":")
				if isList {
					fmt.Print("[]")
				}
				g.Types[typeName] = tokens.GQLTypeToGolangType(typeName)
				fmt.Print(fieldTypeName, "\n\t")
			}
		}
		fmt.Println()
	}
}

func (g *Generator) HandleEnum(req interface{}) {
	z := req.(map[string]interface{})
	fmt.Print(z["kind"], " : ", z["name"], "\n")
	enumList := []string{}
	enumValues := tm(req)["enumValues"]
	if enumValues != nil {
		for _, val := range enumValues.([]interface{}) {
			enumList = append(enumList, tm(val)["name"].(string))
		}
	}
	fmt.Print("Values: [", strings.Join(enumList, ", "), "]\n\n")
	return
}
