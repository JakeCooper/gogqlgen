package tokens

var primitiveMap = map[string]string{
	"String":   "string",
	"BigInt":   "int64",
	"Json":     "map[string]interface{}",
	"Boolean":  "bool",
	"Float":    "float64",
	"Int":      "int32",
	"DateTime": "string",
}

func GQLTypeToGolangType(gqlType string) string {
	if v, ok := primitiveMap[gqlType]; ok {
		return v
	}
	return gqlType
}

func CanGenerate(gqlType string) bool {
	_, ok := primitiveMap[gqlType]
	return ok
}
