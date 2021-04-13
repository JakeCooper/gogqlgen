package generator

import (
	"fmt"
	"strings"

	"github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/tokens"
)

type ArgShit struct {
	Name       string
	TypeName   string
	IsRequired bool
}

func (g *Generator) HandleQuery(req interface{}) {
	z := req.(map[string]interface{})
	fields := z["fields"]
	// name := z["name"]
	// Handle Everything Else
	// g.Outs.ClientFile.Write([]byte("import (\n \"context\"\n)\n\n"))
	// g.Outs.ClientFile.Write([]byte("type Client struct {}\n\n"))
	if fields != nil {
		for _, rf := range fields.([]interface{}) {
			field := tm(rf)
			fieldType := tm(field["type"])["ofType"]
			args := field["args"].([]interface{})
			processedArgs := make([]ArgShit, len(args))
			fieldName := field["name"].(string)
			g.Outs.QueryFile.Write([]byte(fmt.Sprintf("type %sRequest struct {\n", strings.Title(fieldName))))
			for i, arg := range args {
				argStruct := tm(arg)
				argName := argStruct["name"]
				argValue := tm(argStruct["type"])["name"]
				isRequired := tm(argStruct["type"])["kind"] == "NON_NULL"
				if argValue == nil {
					argValue = tm(tm(argStruct["type"])["ofType"])["name"]
				}
				processedArgs[i] = ArgShit{
					Name:       argName.(string),
					TypeName:   argValue.(string),
					IsRequired: isRequired,
				}
				g.Outs.QueryFile.Write([]byte("\t"))
				g.Outs.QueryFile.Write([]byte(strings.Title(argName.(string))))
				g.Outs.QueryFile.Write([]byte(" "))
				if !isRequired {
					g.Outs.QueryFile.Write([]byte("*"))
				}
				g.Outs.QueryFile.Write([]byte(tokens.GQLTypeToGolangType(argValue.(string))))
				g.Outs.QueryFile.Write([]byte("\n"))
			}
			// Whatever the response is GQL
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

				if graphql.IsInternal(typeName) || tokens.IsPrimitive(typeName) {
					// Skip internal GraphQL shit
					g.writeClient(fieldName, typeName, isList, processedArgs, !tokens.IsPrimitive(typeName))

					g.Outs.QueryFile.Write([]byte("}\n\n"))
					continue
				}
				g.Outs.QueryFile.Write([]byte("\tGQL "))

				if isList {
					// g.Outs.QueryFile.Write([]byte("[]"))
				}
				// g.Types[typeName] = typeName
				// if !tokens.IsPrimitive(typeName) {
				// 	g.Outs.QueryFile.Write([]byte("*"))
				// }
				g.Outs.QueryFile.Write([]byte(fmt.Sprintf("%sGQL\n", tokens.GQLTypeToGolangType(typeName))))
				g.Outs.QueryFile.Write([]byte("}\n\n"))
				// TODO Generate client (Probably need a response object else gonna fail to serialize shit back)
				// Actually probs not true just do the machinebox gql meme with the resp as (bool, primitive, etc)
				g.writeClient(fieldName, typeName, isList, processedArgs, !tokens.IsPrimitive(typeName))
			}
		}
	}
}

func (g *Generator) writeClient(name string, typeName string, isList bool, args []ArgShit, hasQueryable bool) {
	upperName := strings.Title(name)
	typeValue := tokens.GQLTypeToGolangType(typeName)
	if isList {
		typeValue = "[]" + typeValue
	} else {
		typeValue = "*" + typeValue
	}
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("func (c *Client) %s (ctx context.Context, req *%sRequest) (%s, error) {\n", upperName, upperName, typeValue)))
	if hasQueryable {
		g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\tgen, err := c.asGQL(ctx, req.GQL)\n")))
		g.Outs.ClientFile.Write([]byte("\tif err != nil {\n\t\treturn nil, err\n\t}\n"))
	}
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\tgqlreq := graphql.NewRequest(fmt.Sprintf(`\n\t\tquery %s {\n", getQueryParamString(args))))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\t\t\t%s%s%s\n", name, getInnerQueryString(args), getQuerySelection(hasQueryable))))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\t\t}`%s))\n", getGen(hasQueryable))))
	for _, arg := range args {
		g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\tgqlreq.Var(\"%s\", req.%s)\n", arg.Name, strings.Title(arg.Name))))
	}
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\tvar resp struct {\n")))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\t\t%s %s%s `json:\"%s\"`\n\t}\n", strings.Title(name), arrayChar(isList), tokens.GQLTypeToGolangType(typeName), name)))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\tif err := c.gql.Run(ctx, gqlreq, &resp); err != nil {\n")))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\t\treturn nil, errors.New(\"%s request failed\")\n", strings.Title(name))))
	g.Outs.ClientFile.Write([]byte("\t}\n"))
	g.Outs.ClientFile.Write([]byte(fmt.Sprintf("\treturn resp.%s, nil\n", strings.Title(name))))
	g.Outs.ClientFile.Write([]byte("}\n\n"))
}

func getGen(hasQueryable bool) string {
	if hasQueryable {
		return ", *gen"
	}
	return ""
}

func arrayChar(isArray bool) string {
	if isArray {
		return "[]"
	}
	return "*"
}

func getQuerySelection(hasQuery bool) string {
	if hasQuery {
		return " {\n\t\t\t\t%s\n\t\t\t}"
	}
	return ""
}

func getQueryParamString(args []ArgShit) string {
	str := "("
	kablarg := []string{}
	for _, arg := range args {
		interstring := "$" + arg.Name + ": " + arg.TypeName
		if arg.IsRequired {
			interstring += "!"
		}
		kablarg = append(kablarg, interstring)
	}
	if len(kablarg) == 0 {
		return ""
	}
	str += strings.Join(kablarg, ", ")
	str += ")"
	return str
}

func getInnerQueryString(args []ArgShit) string {
	str := "("
	kablarg := []string{}
	for _, arg := range args {
		interstring := arg.Name + ": $" + arg.Name
		kablarg = append(kablarg, interstring)
	}
	if len(kablarg) == 0 {
		return ""
	}
	str += strings.Join(kablarg, ", ")
	str += ")"
	return str
}

func generateSelectionIfNeeded() string {
	// If it's a primitive type, no GQL and no selection required
	// Otherwise use janky serializer
	return ""
}
