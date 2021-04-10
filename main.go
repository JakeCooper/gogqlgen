package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jakecooper/gogqlgen/internal/introspect"
)

func main() {
	_, err := ioutil.ReadFile("./operations.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := introspect.RawSchema("http://localhost:8082")
	if err != nil {
		panic(err)
	}
	types := schema["__schema"].(map[string]interface{})["types"].([]interface{})
	for _, v := range types {
		z := v.(map[string]interface{})
		kind := z["kind"]
		name := z["name"]

		fmt.Print(kind, ": ", name, "\n")
		f := getFields(z)
		if f != nil {
			fmt.Print("Fields: ")
			for _, rf := range f.([]interface{}) {
				field := rf.(map[string]interface{})
				fieldType := field["type"].(map[string]interface{})["ofType"]
				if fieldType != nil {
					fieldTypeName := fieldType.(map[string]interface{})["name"]
					fmt.Print(field["name"], ":", fieldTypeName, "\n\t")
				}
			}
			fmt.Println()
		}
	}
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
