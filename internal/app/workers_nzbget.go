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
func (j *NzbgetProcess) Work(ctx context.Context, job *minion.Job[*NzbgetProcess]) error {
	l := app.Log.Named("workers.nzbget")
	p := job.Args.Payload

	dir := p.Dir
	if p.FinalDir != "" {
		dir = p.FinalDir
	}
	dir = app.Config.DirectoriesNzbget + dir

	download, err := app.DB.DownloadByHash(p.Id)
	if err != nil {
		return fae.Wrap(err, "getting download by hash")
	}
	if download == nil {
		return fae.Errorf("download not found: %s", p.Id)
	}
	if download.Status == "reviewing" {
		return fae.Errorf("download reviewing: %s", p.Id)
	}
	if download.Status == "done" {
		return fae.Errorf("download done: %s", p.Id)
	}

	if download.Medium == nil {
		return fae.Errorf("download has no medium: %s", p.Id)
	}

	if p.Status != "SUCCESS" {
		download.Status = "reviewing"
		if err := app.DB.Download.Save(download); err != nil {
			return fae.Wrap(err, "saving download")
		}

		return nil
	}

	l.Debugf("nzbget process: %s: %s %s", p.Id, download.Medium.Title, download.Medium.Display)

	// list all files in dir
	files, err := j.GetFiles(dir)
	if err != nil {
		return fae.Wrap(err, "getting files")
	}
	if len(files) == 0 {
		return fae.Errorf("no files found: %s", dir)
	}
	if len(files) > 1 {
		return fae.Errorf("multiple files found: %s", dir)
	}

	file := files[0]
	ext := filepath.Ext(file)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	dest, err := Destination(download.Medium)
	if err != nil {
		return fae.Wrap(err, "getting destination")
	}

	destination := filepath.Join(app.Config.DirectoriesCompleted, fmt.Sprintf("%s.%s", dest, ext))
	source := filepath.Join(dir, file)
	l.Debugf("nzbget process: %s: %s => %s", p.Id, source, destination)

	if !app.Config.Production {
		l.Debugf("skipping move in dev mode")
		return nil
	}

	if err := FileLink(source, destination, true); err != nil {
		return fae.Wrap(err, "linking file")
	}

	if err := updateMedium(download.Medium, []string{dest}); err != nil {
		return fae.Wrap(err, "updating medium")
	}

	if err := app.Plex.RefreshLibraryPath(filepath.Dir(destination)); err != nil {
		return fae.Wrap(err, "refreshing plex library")
	}

	download.Status = "done"
	if err := app.DB.Download.Save(download); err != nil {
		return fae.Wrap(err, "saving download")
	}

	notifier.Success("Downloads::Completed", fmt.Sprintf("%s %s", download.Medium.Title, download.Medium.Display))
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
