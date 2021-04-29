package nagios

import "net/url"

type Query struct {
	Endpoint string
	URLQuery url.Values
}

func (q Query) SetNonEmpty(key, value string) {
	if len(value) > 0 && value != "0" {
		q.URLQuery.Set(key, value)
	}
}

type QueryBuilder interface {
	Build() Query
}
