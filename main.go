package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jakecooper/gogqlgen/internal/introspect"
	"github.com/jakecooper/gogqlgen/internal/tokens"
)

func main() {
	_, err := ioutil.ReadFile("./operations.graphql")
	if err != nil {
		panic(err)
	}
	url := flag.String("url", "", "URL of GraphQL API")
	flag.Parse()
	if url == nil || *url == "" {
		panic(errors.New("URL must be provided!"))
	}
	schema, err := introspect.RawSchema(*url)
	if err != nil {
		panic(err)
	}
	types := schema["__schema"].(map[string]interface{})["types"].([]interface{})
	typeMap := make(map[string]string)
	for _, v := range types {
		z := v.(map[string]interface{})
		kind := z["kind"]
		name := z["name"]

		// Skip internal graphql shit
		if name == "__Type" || name == "__Schema" {
			continue
		}

		fmt.Print(kind, ": ", name, "\n")
		f := getFields(z)

		// TODO Probs a switch statement
		// Handle Enum
		if kind == "ENUM" {
			enumList := []string{}
			enumValues := v.(map[string]interface{})["enumValues"]
			if enumValues != nil {
				for _, val := range enumValues.([]interface{}) {
					enumList = append(enumList, val.(map[string]interface{})["name"].(string))
				}
			}
			fmt.Print("Values: [", strings.Join(enumList, ", "), "]\n\n")
			continue
		}

		// Handle Everything Else
		if f != nil {
			fmt.Print("Fields: ")
			for _, rf := range f.([]interface{}) {
				field := rf.(map[string]interface{})
				fieldType := field["type"].(map[string]interface{})["ofType"]
				if fieldType != nil {
					fieldTypeName := fieldType.(map[string]interface{})["name"]
					fieldKind := fieldType.(map[string]interface{})["kind"]
					isList := false
					if fieldKind != nil {
						// InternalType
						// TODO Probably recursion IDK
						ofType := fieldType.(map[string]interface{})["ofType"]
						isList = fieldType.(map[string]interface{})["kind"] == "LIST"
						if ofType != nil {
							ofType = ofType.(map[string]interface{})["ofType"]
							if ofType != nil {
								fieldTypeName = ofType.(map[string]interface{})["name"]
							}
						}
					}
					typeName := fieldTypeName.(string)
					// TODO just check if the prefix is __

					if strings.Contains(typeName, "__") {
						// Skip internal GraphQL shit
						continue
					}
					fmt.Print(field["name"], ":")
					if isList {
						fmt.Print("[]")
					}
					typeMap[typeName] = tokens.GQLTypeToGolangType(typeName)
					fmt.Print(fieldTypeName, "\n\t")
				}
			}
			fmt.Println()
		}
	}

	b, err := json.MarshalIndent(typeMap, "", "  ")
	if err != nil {
		panic("FAILED TO RANDY MARSHAL")
	}
	fmt.Println("List o Types", string(b))
}

func getFields(req interface{}) interface{} {
	return req.(map[string]interface{})["fields"]
}

func tm(req interface{}) map[string]interface{} {
	return req.(map[string]interface{})
}

func toStruct(mp map[string]interface{}, resp interface{}) error {
	b, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &resp)
}

type GQLType struct {
	Kind       string        `json:"kind"`
	Name       string        `json:"name"`
	Decription string        `json:"description"`
	Fields     []interface{} `json:"fields"`
}
