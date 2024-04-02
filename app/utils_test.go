package app

import "testing"

func TestPath(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"easy", "easy"},
		{"don't", "dont"},
		{"don't worry", "dont worry"},
		{"I'm sorry", "im sorry"},
		{"I'm sorry!", "im sorry"},
		{"You're welcome", "youre welcome"},
		{"You're welcome!", "youre welcome"},
		{" Leading Spaces", "leading spaces"},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if got := path(tt.title); got != tt.want {
				t.Errorf("path() = %v, want %v", got, tt.want)
			}
		})
	}
}
