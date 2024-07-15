package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestPathParts(t *testing.T) {
	tests := []struct {
		title string
		kind  string
		name  string
		file  string
		ext   string
	}{
		// {
		// 	"/mnt/media/anime/my dress up darling/my dress-up darling - 01x001 - someone who lives in the exact opposite world as me.mkv",
		// 	"anime",
		// 	"my dress up darling",
		// 	"my dress-up darling - 01x001 - someone who lives in the exact opposite world as me",
		// 	"mkv",
		// },
		{
			"/mnt/media/anime/arifureta from commonplace to worlds strongest/arifureta - from commonplace to worlds strongest - 00x001 - omnibus the great orcus labyrinth episode 5.5.mkv",
			"anime",
			"arifureta from commonplace to worlds strongest",
			"arifureta - from commonplace to worlds strongest - 00x001 - omnibus the great orcus labyrinth episode 5.5",
			"mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			kind, name, file, ext, err := pathParts(tt.title)
			assert.NoError(t, err)
			assert.Equal(t, kind, tt.kind)
			assert.Equal(t, name, tt.name)
			assert.Equal(t, file, tt.file)
			assert.Equal(t, ext, tt.ext)
		})
	}
}
