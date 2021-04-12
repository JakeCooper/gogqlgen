package generator

import (
	"fmt"
	"strings"

	"github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/tokens"
)

func (g *Generator) HandleObject(req interface{}) {
	z := req.(map[string]interface{})
	name := z["name"]
	fields := z["fields"]

	// TODO also generate a <Object>GQL struct { field: bool, field: WhateverGQL }

	g.Outs.ObjectFile.Write([]byte(fmt.Sprintf("type %s struct {\n", name.(string))))
	g.Outs.GQLFile.Write([]byte(fmt.Sprintf("type %sGQL struct {\n", name.(string))))
	// Handle Everything Else
	if fields != nil {
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
				fieldName := field["name"].(string)
				g.Outs.ObjectFile.Write([]byte(fmt.Sprintf("\t%s ", strings.Title(fieldName))))
				g.Outs.GQLFile.Write([]byte(fmt.Sprintf("\t%s ", strings.Title(fieldName))))

				if isList {
					g.Outs.ObjectFile.Write([]byte("[]"))
				}
				g.Types[typeName] = typeName
				if !tokens.IsPrimitive(typeName) {
					g.Outs.ObjectFile.Write([]byte("*"))
					g.Outs.GQLFile.Write([]byte(fmt.Sprintf("*%sGQL `json:\"%s\"`\n", typeName, fieldName)))
				} else {
					g.Outs.GQLFile.Write([]byte(fmt.Sprintf(fmt.Sprintf("bool `json:\"%s\"`\n", fieldName))))
				}
				g.Outs.ObjectFile.Write([]byte(fmt.Sprintf("%s `json:\"%s\"`\n", tokens.GQLTypeToGolangType(typeName), fieldName)))
			}
		}
	}
	g.Outs.ObjectFile.Write([]byte("}\n\n"))
	g.Outs.GQLFile.Write([]byte("}\n\n"))
}
