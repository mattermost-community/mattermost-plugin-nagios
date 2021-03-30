package nagios

import (
	"net/url"
	"testing"
)

func TestQuery_SetNonEmpty(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name     string
		URLQuery url.Values
		args     args
	}{
		{
			name:     "empty value",
			URLQuery: make(url.Values),
			args: args{
				key: "a",
			},
		},
		{
			name:     "non-empty value",
			URLQuery: make(url.Values),
			args: args{
				key:   "a",
				value: "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Query{URLQuery: tt.URLQuery}

			q.SetNonEmpty(tt.args.key, tt.args.value)

			if v, ok := q.URLQuery[tt.args.key]; ok {
				if len(tt.args.value) == 0 {
					t.Errorf("URLQuery[%s] should be empty", tt.args.key)
				} else if v[0] != tt.args.value {
					t.Errorf("got %v, want %v", v[0], tt.args.key)
				}
			} else if len(tt.args.value) > 0 {
				t.Errorf("URLQuery[%s] should not be empty", tt.args.key)
			}
		})
	}
}
