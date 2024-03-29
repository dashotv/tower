package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

	"github.com/dashotv/fae"
)

func imageDownload(source, destination string) error {
	base := filepath.Dir(destination)
	if _, err := os.Stat(base); os.IsNotExist(err) {
		if err := os.MkdirAll(base, 0755); err != nil {
			return fae.Wrap(err, "creating directory")
		}
	}

	out, err := os.Create(destination)
	if err != nil {
		return fae.Wrap(err, "creating file")
	}
	defer out.Close()

	resp, err := http.Get(source)
	if err != nil {
		return fae.Wrap(err, "getting image")
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fae.Wrap(err, "copying image")
	}

	err = os.Chown(destination, 1001, 1001)
	if err != nil {
		return fae.Wrap(err, "chowning image")
	}

	return nil
}

func imageResize(source, destination string, width, height int) error {
	img, err := imgio.Open(source)
	if err != nil {
		return fae.Wrap(err, "opening image")
	}

	resized := transform.Resize(img, width, height, transform.Lanczos)
	if err := imgio.Save(destination, resized, imgio.JPEGEncoder(80)); err != nil {
		return fae.Wrap(err, "saving image")
	}

	err = os.Chown(destination, 1001, 1001)
	if err != nil {
		return fae.Wrap(err, "chowning image")
	}

	return nil
}
