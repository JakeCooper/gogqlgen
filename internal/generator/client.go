package generator

import "fmt"

// TODO generate the whole fkn thing with machineboxgql
func GenerateClient(name string) string {
	return fmt.Sprintf("func (c *Client) %s (ctx context.Context, req %sRequest) *%s {}", name, name, name)
}
