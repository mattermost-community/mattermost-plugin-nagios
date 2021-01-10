package main

import "testing"

func Test_emoji(t *testing.T) {
	tests := []struct {
		name  string
		state string
		want  string
	}{
		{
			name:  "empty",
			state: "",
			want:  questionEmoji,
		},
		{
			name:  "unknown",
			state: "testing",
			want:  questionEmoji,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := emoji(tt.state); got != tt.want {
				t.Errorf("emoji() = %v, want %v", got, tt.want)
			}
		})
	}
}
