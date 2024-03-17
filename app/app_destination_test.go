package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayout(t *testing.T) {
	cases := []struct {
		id string
		l  string
	}{
		{id: "65a0f742175ec2916ae434b8", l: "anime/shangri la frontier/shangri la frontier"},
		{id: "65a0f745175ec2916ae434d6", l: "anime/shangri la frontier/shangri la frontier - 01x23 #023 - Bird with Rabbits vs. Skeletal Choir"},
		{id: "58bdf1df6b696d7139000000", l: "movies4k/arrival/arrival"},
		{id: "655a5b473359bb31b6f4932a", l: "tv/doctor who 2005/doctor who 2005 - 14x01 - TBA"},
	}

	err := appSetup()
	require.NoError(t, err)

	for _, tt := range cases {
		t.Run(tt.l, func(t *testing.T) {
			m := &Medium{}
			err := app.DB.Medium.Find(tt.id, m)
			require.NoError(t, err)

			l, err := Destination(m)
			assert.NoError(t, err)
			assert.Equal(t, tt.l, l)
		})
	}
}
