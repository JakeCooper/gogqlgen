package generator

func GenerateClient(name string) string {
	return "func (c *Client) %s (ctx context.Context, req %s) *%s {}"
}
