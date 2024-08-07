package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

type NzbgetProcess struct {
	minion.WorkerDefaults[*NzbgetProcess]
	Payload *NzbgetPayload `json:"payload"`
}

func (j *NzbgetProcess) Kind() string { return "nzbget_process" }
func (j *NzbgetProcess) Work(ctx context.Context, job *minion.Job[*NzbgetProcess]) (err error) {
	a := ContextApp(ctx)
	l := a.Log.Named("workers.nzbget")
	p := job.Args.Payload

	dir := p.Dir
	if p.FinalDir != "" {
		dir = p.FinalDir
	}
	dir = a.Config.DirectoriesNzbget + dir

	download, err := a.DB.DownloadByHash(p.ID)
	if err != nil {
		return fae.Wrap(err, "getting download by hash")
	}
	if download == nil {
		return fae.Errorf("download not found: %s", p.ID)
	}
	if download.Status == "reviewing" {
		return fae.Errorf("download reviewing: %s", p.ID)
	}
	if download.Status == "done" {
		return fae.Errorf("download done: %s", p.ID)
	}

	if download.Medium == nil {
		return fae.Errorf("download has no medium: %s", p.ID)
	}

	if p.Status != "SUCCESS" {
		download.Status = "reviewing"
		if err := a.DB.Download.Save(download); err != nil {
			return fae.Wrap(err, "saving download")
		}

		return nil
	}

	defer func() {
		if err != nil {
			download.Status = "reviewing"
			if err := a.DB.Download.Save(download); err != nil {
				l.Error("saving download", "error", err)
			}
		}
	}()

	l.Debugf("nzbget process: %s: %s %s", p.ID, download.Title, download.Display)

	// list all files in dir
	files, err := j.GetFiles(dir)
	if err != nil {
		return fae.Wrap(err, "getting files")
	}
	if len(files) == 0 {
		return fae.Errorf("no files found: %s", dir)
	}
	if len(files) > 1 {
		return fae.Errorf("multiple files found: %s: %+v", dir, files)
	}

	// TODO: use Mover
	file := files[0]
	ext := filepath.Ext(file)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	dest, err := a.Destinator.Destination(download.Kind, download.Medium)
	if err != nil {
		return fae.Wrap(err, "getting destination")
	}

	destination := fmt.Sprintf("%s.%s", dest, ext)
	source := filepath.Join(dir, file)
	l.Debugf("nzbget process: %s: %s => %s", p.ID, source, destination)

	if !a.Config.Production {
		l.Debugf("skipping move in dev mode")
		return nil
	}

	if err := FileLink(source, destination, true); err != nil {
		return fae.Wrap(err, "linking file")
	}

	if err := a.updateMedia([]*MoverFile{{Medium: download.Medium, Source: source, Destination: destination}}); err != nil {
		return fae.Wrap(err, "updating medium")
	}

	if err := a.Plex.RefreshLibraryPath(filepath.Dir(destination)); err != nil {
		return fae.Wrap(err, "refreshing plex library")
	}

	download.Status = "done"
	if err := a.DB.Download.Save(download); err != nil {
		return fae.Wrap(err, "saving download")
	}

	notifier.Success("Downloads::Completed", fmt.Sprintf("%s - %s", download.Title, download.Display))
	return nil
}

func (j *NzbgetProcess) GetFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fae.Wrap(err, "reading dir")
	}

	var list []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		list = append(list, f.Name())
	}

	list = lo.Filter(list, func(s string, i int) bool {
		return shouldDownloadFile(s)
	})

	// TODO: handle multiple files? subtitles? also, figure nzbs of full seasons

	return list, nil
}
