package nagios

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Query struct{ url.Values }

type QueryBuilder interface {
	Build() Query
}

type Client struct {
	c *http.Client
	u *url.URL
}

func cloneURLUpToPath(u *url.URL) *url.URL {
	return &url.URL{
		Scheme:      u.Scheme,
		Opaque:      u.Opaque,
		User:        u.User,
		Host:        u.Host,
		Path:        u.Path,
		RawPath:     u.RawPath,
		ForceQuery:  true,
		RawQuery:    "",
		Fragment:    "",
		RawFragment: "",
	}
}

func (c Client) Query(q QueryBuilder, v interface{}) error {
	u := cloneURLUpToPath(c.u)

	u.RawQuery = q.Build().Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %w", err)
	}

	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("Do: %w", err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return fmt.Errorf("Decode: %w", err)
	}

	return nil
}

func NewClient(client *http.Client, address string) (*Client, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}

	return &Client{
		c: client,
		u: cloneURLUpToPath(u),
	}, nil
}
