package generator

import (
	"fmt"
	"strings"

	"github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/tokens"
)

func (g *Generator) HandleQuery(req interface{}) {
	z := req.(map[string]interface{})
	fields := z["fields"]
	// name := z["name"]
	// Handle Everything Else
	if fields != nil {
		for _, rf := range fields.([]interface{}) {
			field := tm(rf)
			fmt.Println(field["name"])
			fieldType := tm(field["type"])["ofType"]
			args := field["args"].([]interface{})
			fieldName := field["name"].(string)
			g.Outs.QueryFile.Write([]byte(fmt.Sprintf("type %sRequest struct {\n", strings.Title(fieldName))))
			// g.Outfile.Write([]byte(fmt.Sprintf("type %sRequest struct {\n", name.(string))))
			for _, arg := range args {
				argStruct := tm(arg)
				argName := argStruct["name"]
				argValue := tm(argStruct["type"])["name"]
				isRequired := tm(argStruct["type"])["kind"] == "NON_NULL"
				if argValue == nil {
					argValue = tm(tm(argStruct["type"])["ofType"])["name"]
				}
				g.Outs.QueryFile.Write([]byte("\t"))
				g.Outs.QueryFile.Write([]byte(strings.Title(argName.(string))))
				g.Outs.QueryFile.Write([]byte(" "))
				if !isRequired {
					g.Outs.QueryFile.Write([]byte("*"))
				}
				g.Outs.QueryFile.Write([]byte(tokens.GQLTypeToGolangType(argValue.(string))))
				g.Outs.QueryFile.Write([]byte("\n"))
				fmt.Println(argName, argValue, isRequired)
			}
			// Whatever the response is GQL
			if fieldType != nil {
				fieldTypeName := tm(fieldType)["name"]
				fieldKind := tm(fieldType)["kind"]
				// isList := false
				if fieldKind != nil {
					// InternalType
					// TODO Probably recursion IDK
					ofType := tm(fieldType)["ofType"]
					// isList = tm(fieldType)["kind"] == "LIST"
					if ofType != nil {
						ofType = tm(ofType)["ofType"]
						if ofType != nil {
							fieldTypeName = tm(ofType)["name"]
						}
					}
				}
				typeName := fieldTypeName.(string)

				if graphql.IsInternal(typeName) || tokens.IsPrimitive(typeName) {
					// Skip internal GraphQL shit
					g.Outs.QueryFile.Write([]byte("}\n\n"))
					continue
				}
				g.Outs.QueryFile.Write([]byte("\tGQL "))

				// if isList {
				// 	g.Outs.QueryFile.Write([]byte("[]"))
				// }
				// g.Types[typeName] = typeName
				// if !tokens.IsPrimitive(typeName) {
				// 	g.Outs.QueryFile.Write([]byte("*"))
				// }
				g.Outs.QueryFile.Write([]byte(fmt.Sprintf("%sGQL\n", tokens.GQLTypeToGolangType(typeName))))
				g.Outs.QueryFile.Write([]byte("}\n\n"))
			}
		}
	}
}
