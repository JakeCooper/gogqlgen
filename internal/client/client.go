package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/machinebox/graphql"
)

type Client struct {
	headers map[string]string
	gql     *graphql.Client
}

func NewClient(endpoint string, opts ...graphql.ClientOption) *Client {
	gql := graphql.NewClient(endpoint, opts...)
	return &Client{
		gql:     gql,
		headers: make(map[string]string),
	}
}

func (c *Client) WithHeader(key string, value string) *Client {
	c.headers[key] = value
	return c
}

func (c *Client) doRequest(ctx context.Context, req *graphql.Request, res interface{}) error {
	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	c.headers = make(map[string]string)
	if err := c.gql.Run(ctx, req, &res); err != nil {
		return err
	}
	return nil
}

func (c *Client) asGQL(ctx context.Context, req interface{}) (*string, error) {
	mp := make(map[string]interface{})
	bytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &mp)
	if err != nil {
		return nil, err
	}
	fields := []string{}
	for k, i := range mp {
		// GQL Selection
		switch i.(type) {
		case bool:
			// GQL Selection
			fields = append(fields, k)
		case map[string]interface{}:
			// Nested GQL/Struct
			nested, err := c.asGQL(ctx, i)
			if err != nil {
				return nil, err
			}
			fields = append(fields, fmt.Sprintf("%s {\n%s\n}", k, *nested))
		default:
			return nil, errors.New("Unsupported Type! Cannot generate GQL")
		}
	}
	q := strings.Join(fields, "\n")
	return &q, nil
}
