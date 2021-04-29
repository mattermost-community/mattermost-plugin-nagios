package nagios

import "testing"

func Test_buildOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []string
		want    string
	}{
		{
			name:    "none",
			options: nil,
			want:    "",
		},
		{
			name:    "some",
			options: []string{"alice", "bob"},
			want:    "alice bob",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildOptions(tt.options); got != tt.want {
				t.Errorf("buildOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatOptions_String(t *testing.T) {
	tests := []struct {
		name          string
		formatOptions FormatOptions
		want          string
	}{
		{
			name:          "none",
			formatOptions: FormatOptions{},
			want:          "",
		},
		// TODO(amwolff): add test cases.
		{
			name: "all",
			formatOptions: FormatOptions{
				Whitespace: true,
				Enumerate:  true,
				Bitmask:    true,
				Duration:   true,
			},
			want: "whitespace enumerate bitmask duration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.formatOptions.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
