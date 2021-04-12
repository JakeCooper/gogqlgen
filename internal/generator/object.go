package generator

import (
	"fmt"

	"github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/tokens"
)

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
