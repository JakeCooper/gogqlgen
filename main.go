package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/jakecooper/gogqlgen/internal/generator"
	gql "github.com/jakecooper/gogqlgen/internal/graphql"
	"github.com/jakecooper/gogqlgen/internal/introspect"
	"github.com/jakecooper/gogqlgen/internal/mapper"
)

var tm = mapper.ToMap

func main() {
	url := flag.String("url", "", "URL of GraphQL API")
	flag.Parse()
	if url == nil || *url == "" {
		panic(errors.New("URL must be provided!"))
	}

	g := generator.New("gen", "generated")

	schema, err := introspect.RawSchema(*url)
	if err != nil {
		panic(err)
	}
	types := schema["__schema"].(map[string]interface{})["types"].([]interface{})
	for _, v := range types {
		z := tm(v)
		kind := z["kind"]
		name := z["name"]

		// Skip internal graphql shit
		if name == "__Type" || name == "__Schema" {
			continue
		}

		if gql.IsInternal(name.(string)) {
			// Skip internal GraphQL shit
			continue
		}

		// Handle Queries, mutations, etc
		switch name {
		case "Mutation":
			g.HandleMutation(z)
			continue
		case "Query":
			g.HandleQuery(z)
			continue
		default:
			// Continue onward
		}

		// Handle primitives
		switch kind {
		case "INPUT_OBJECT":
			g.HandleInputObject(z)
		case "SCALAR":
			g.HandleScalar(z)
		case "ENUM":
			g.HandleEnum(z)
		default:
			g.HandleObject(z)
		}
	}

	fmt.Println("Generation Successfully! (Probably)")
}
