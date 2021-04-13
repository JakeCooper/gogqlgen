package generator

import "fmt"

// TODO generate the whole fkn thing with machineboxgql
func (g *Generator) GenerateClient(name string) string {
	return fmt.Sprintf("func (c *Client) %s (ctx context.Context, req *%sRequest) *%s {}\n\n", name, name, name)
}
