package nagios

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func mustParseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(fmt.Sprintf("Parse: %v", err))
	}
	return u
}

func Test_cloneURLToPath(t *testing.T) {
	tests := []struct {
		name string
		u    *url.URL
		want *url.URL
	}{
		{
			name: "u and want are the same",
			u:    mustParseURL("https://github.com"),
			want: mustParseURL("https://github.com"),
		},
		{
			name: "u has a path",
			u:    mustParseURL("https://github.com/ulumuri/go-nagios"),
			want: mustParseURL("https://github.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cloneURLToPath(tt.u); got.String() != tt.want.String() {
				t.Errorf("cloneURLToPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPathLayout(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "Archive JSON CGI",
			endpoint: "archivejson.cgi",
			want:     "/nagios/cgi-bin/archivejson.cgi",
		},
		{
			name:     "Object JSON CGI",
			endpoint: "objectjson.cgi",
			want:     "/nagios/cgi-bin/objectjson.cgi",
		},
		{
			name:     "Status JSON CGI",
			endpoint: "statusjson.cgi",
			want:     "/nagios/cgi-bin/statusjson.cgi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPathLayout(tt.endpoint); got != tt.want {
				t.Errorf("getPathLayout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		client  *http.Client
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "valid address",
			args: args{
				client:  &http.Client{},
				address: "https://nagios.fedoraproject.org/nagios/",
			},
			want: &Client{
				c: &http.Client{},
				u: mustParseURL("https://nagios.fedoraproject.org"),
			},
			wantErr: false,
		},
		{
			name: "invalid address",
			args: args{
				client:  &http.Client{},
				address: "\n",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.client, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const testEndpoint = "TestJSON.CGI"

type mockQueryBuilder struct{}

func (m mockQueryBuilder) Build() Query {
	return Query{
		Endpoint: testEndpoint,
		URLQuery: url.Values{"foo": []string{"bar"}},
	}
}

func TestClient_Query(t *testing.T) {
	type response struct {
		ExampleField string `json:"examplefield"`
	}

	want := response{ExampleField: "examplevalue"}

	h := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query()["foo"][0] != "bar" {
			t.Errorf("Request URL Query does not contain disared values")
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
		}
	}

	mux := http.NewServeMux()
	mux.Handle(getPathLayout(testEndpoint), http.HandlerFunc(h))

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	c, err := NewClient(server.Client(), server.URL)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var (
		m   mockQueryBuilder
		got response
	)

	if err := c.Query(m, &got); err != nil {
		t.Fatalf("Query: %v", err)
	}

	if got.ExampleField != want.ExampleField {
		t.Errorf("got = %v, want %v", got, want)
	}
}
