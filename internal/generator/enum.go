package generator

import (
	"fmt"
	"strings"
)

func (g *Generator) HandleEnum(req interface{}) {
	z := req.(map[string]interface{})
	fmt.Print(z["kind"], " : ", z["name"], "\n")
	enumList := []string{}
	enumValues := tm(req)["enumValues"]
	if enumValues != nil {
		for _, val := range enumValues.([]interface{}) {
			enumList = append(enumList, tm(val)["name"].(string))
		}
	}
	fmt.Print("Values: [", strings.Join(enumList, ", "), "]\n\n")
	return
}
