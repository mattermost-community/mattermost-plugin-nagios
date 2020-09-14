package nagios

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Query struct {
	Endpoint string
	URLQuery url.Values
}

type QueryBuilder interface {
	Build() Query
}

// Client represents a Nagios Core JSON CGIs client.
type Client struct {
	c *http.Client
	u *url.URL
}

func cloneURLToPath(u *url.URL) *url.URL {
	return &url.URL{
		Scheme: u.Scheme,
		Opaque: u.Opaque,
		User:   u.User,
		Host:   u.Host,
	}
}

func getPathLayout(endpoint string) string {
	return fmt.Sprintf("/nagios/cgi-bin/%s", endpoint)
}

// Query builds query using QueryBuilder implementation, queries Nagios Core
// instance and stores the response in the compatible value pointed to by v.
func (c Client) Query(b QueryBuilder, v interface{}) error {
	u := cloneURLToPath(c.u)

	q := b.Build()

	u.Path = getPathLayout(q.Endpoint)
	u.RawQuery = q.URLQuery.Encode()

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

// NewClient returns initialized Client and any error encountered.
func NewClient(client *http.Client, address string) (*Client, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}

	return &Client{
		c: client,
		u: cloneURLToPath(u),
	}, nil
}
