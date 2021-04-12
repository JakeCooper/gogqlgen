package generator

import (
	"fmt"
	"strings"
)

func (g *Generator) HandleEnum(req interface{}) {
	z := req.(map[string]interface{})
	enumName := z["name"]

	// Generate Enum Type
	g.Outs.ObjectFile.Write([]byte("// Enum Type\n"))
	g.Outs.ObjectFile.Write([]byte(fmt.Sprintf("type %s string\n", enumName)))

	g.Outs.GQLFile.Write([]byte(fmt.Sprintf("type %sGQL string\n", enumName)))

	// Generate enum
	g.Outs.ObjectFile.Write([]byte("// ENUM Values\n"))
	enumValues := tm(req)["enumValues"]
	g.Outs.ObjectFile.Write([]byte("const (\n"))
	g.Outs.GQLFile.Write([]byte("const (\n"))

	if enumValues != nil {
		for _, val := range enumValues.([]interface{}) {
			valueName := tm(val)["name"].(string)
			g.Outs.ObjectFile.Write([]byte(fmt.Sprintf("\t%s_%s %s = \"%s\"\n", enumName, strings.ToUpper(valueName), enumName, valueName)))
			g.Outs.GQLFile.Write([]byte(fmt.Sprintf("\t%s_%s_GQL %sGQL = \"%s\"\n", enumName, strings.ToUpper(valueName), enumName, valueName)))
		}
	}
	g.Outs.ObjectFile.Write([]byte(")\n\n"))
	g.Outs.GQLFile.Write([]byte(")\n\n"))
	return
}
