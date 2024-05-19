package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type MediumImage struct {
	minion.WorkerDefaults[*MediumImage]
	ID    string
	Type  string
	Path  string
	Ratio float32
}

func (j *MediumImage) Kind() string { return "MediumImage" }
func (j *MediumImage) Work(ctx context.Context, job *minion.Job[*MediumImage]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("no app in context")
	}

	id := job.Args.ID
	t := job.Args.Type
	remote := job.Args.Path
	ratio := job.Args.Ratio

	a.Log.Named("medium_image").Infof("type: %s, remote: %v", t, remote)

	medium, err := app.DB.Medium.Get(id, &Medium{})
	if err != nil {
		return fae.Wrap(err, "getting medium")
	}

	if err := mediumImage(medium, t, remote, ratio); err != nil {
		return fae.Wrap(err, "medium image")
	}

	if err := a.DB.Medium.Save(medium); err != nil {
		return fae.Wrap(err, "saving series")
	}

	return nil
}

func mediumImageID(id string, t string, remote string, ratio float32) error {
	medium, err := app.DB.Medium.Get(id, &Medium{})
	if err != nil {
		return fae.Wrap(err, "getting medium")
	}
	return mediumImage(medium, t, remote, ratio)
}

func mediumImage(medium *Medium, t string, remote string, ratio float32) error {
	extension := filepath.Ext(remote)
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}
	local := fmt.Sprintf("series-%s/%s", medium.ID.Hex(), t)
	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)

	var img *Path
	switch t {
	case "cover":
		img = medium.GetCover()
	case "background":
		img = medium.GetBackground()
	}

	if img != nil && img.Remote == remote {
		return nil
	}

	if img == nil {
		img = &Path{}
		medium.Paths = append(medium.Paths, img)
	}

	img.Type = primitive.Symbol(t)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	if err := imageDownload(remote, dest); err != nil {
		return fae.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return fae.Wrap(err, "resizing image")
	}

	return nil
}

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
