package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath_ParseTag(t *testing.T) {
	tests := []struct {
		name string
		path *Path
		want string
	}{
		{name: "no tag", path: &Path{Local: "series - 01x001 - no tag"}, want: ""},
		{name: "with tags", path: &Path{Local: "series - 01x001 - with tag [subsplease 1080p]"}, want: ""},
		{name: "with tags and dashes", path: &Path{Local: "series - 01x001 - with tag [erai-raws 1080p]"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.path.ParseTag()
			assert.Equal(t, tt.want, tt.path.Tag)
		})
	}
}
