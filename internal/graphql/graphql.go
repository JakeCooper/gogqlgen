package graphql

import "strings"

func IsInternal(name string) bool {
	return strings.Contains(name, "__")
}

type GQLType struct {
	Kind       string        `json:"kind"`
	Name       string        `json:"name"`
	Decription string        `json:"description"`
	Fields     []interface{} `json:"fields"`
}
