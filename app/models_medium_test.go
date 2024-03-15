package app

import (
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func TestConnector_MediumByFile(t *testing.T) {
	c := testConnector()

	list := []string{
		"/mnt/media/donghua/tomb of fallen gods/tomb of fallen gods 01x10.mp4",
		"/mnt/media/donghua/white cat legend/white cat legend - 01x05 #005 - ali 88.mp4",
		"/mnt/media/donghua/throne of seal/throne of seal 02x16.mp4",
		"/mnt/media/donghua/the great ruler/the great ruler 1x25.mp4",
		"/mnt/media/donghua/shrouding the heavens/shrouding the heavens - 01x05 #005 - .mkv",
		// "/mnt/media/donghua/100000 years of refining qi/100000 years of refining qi 01x82.mp4",
	}

	for _, path := range list {
		t.Run(path, func(t *testing.T) {
			f := &File{Path: path}
			m, err := c.MediumByFile(f)
			require.NoError(t, err)
			require.NotNil(t, m)
		})
	}
}
