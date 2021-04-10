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
		z := tm(v)
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
			enumValues := tm(v)["enumValues"]
			if enumValues != nil {
				for _, val := range enumValues.([]interface{}) {
					enumList = append(enumList, tm(val)["name"].(string))
				}
			}
			fmt.Print("Values: [", strings.Join(enumList, ", "), "]\n\n")
			continue
		}

		// Handle Everything Else
		if f != nil {
			fmt.Print("Fields: ")
			for _, rf := range f.([]interface{}) {
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
