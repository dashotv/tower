package app

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/dashotv/tower/internal/plex"
)

func newWalker(db *Connector, logger *zap.SugaredLogger, libs []*plex.PlexLibrary) *Walker {
	return &Walker{
		db:          db,
		logger:      logger,
		Libraries:   libs,
		directories: make(chan string, 10),
		files:       make(chan string, 10),
	}
}

type Walker struct {
	db          *Connector
	logger      *zap.SugaredLogger
	Libraries   []*plex.PlexLibrary
	directories chan string
	files       chan string
}

type counter struct {
	sync.Mutex
	v int
}

func (c *counter) Inc() {
	c.Lock()
	c.v++
	c.Unlock()
}

func (w *Walker) Walk() error {
	c := counter{v: 0}
	start := time.Now()
	defer func() { w.logger.Infow("walk", "duration", time.Since(start), "count", c.v) }()

	eg := new(errgroup.Group)

	eg.Go(func() error {
		defer close(w.directories)

		if err := w.getDirectories(); err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		defer close(w.files)

		for dir := range w.directories {
			if err := w.getFiles(dir); err != nil {
				return err
			}
		}
		return nil
	})

	for i := 0; i < 10; i++ {
		eg.Go(func() error {
			for line := range w.files {
				c.Lock()
				c.v++
				c.Unlock()
				if err := w.updateFile(line); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *Walker) createFile(line string) error {
	file, err := w.db.FileByPath(line)
	if err != nil {
		return err
	}

	if file.ID != primitive.NilObjectID {
		return nil
	}

	if err := w.db.File.Save(file); err != nil {
		return err
	}

	return nil
}

func (w *Walker) updateFile(line string) error {
	file, err := w.db.FileByPath(line)
	if err != nil {
		w.logger.Errorw("file", "error", err)
		return err
	}

	info, err := os.Stat(line)
	if err != nil {
		w.logger.Errorw("stat", "error", err)
		return err
	}

	file.ModifiedAt = info.ModTime().Unix()
	file.Size = info.Size()

	if err := w.db.File.Save(file); err != nil {
		w.logger.Errorw("save", "error", err)
		return err
	}

	return nil
}

func (w *Walker) getDirectories() error {
	out := func(name, line string) {
		// w.logger.Debugf("sending directory: %s", line)
		w.directories <- line
	}
	err := func(name, line string) {
		w.logger.Errorw("shell", "error", line)
	}

	for _, lib := range w.Libraries {
		if lib.Type != "show" && lib.Type != "movie" {
			continue
		}
		for _, loc := range lib.Locations {
			// w.logger.Infow("walking", "library", lib.Title, "path", loc.Path)
			// workers.Shell(w.logger.Named("shell"), loc.Path, "find", loc.Path, "-type", "f", "-exec", "stat", "{}", "+")
			cmd := fmt.Sprintf("find '%s' -maxdepth 1 -mindepth 1 -type d", loc.Path)
			// w.logger.Debugf("running command: %s", cmd)
			status, err := Shell(cmd, ShellOptions{Out: out, Err: err})
			if err != nil {
				return err
			}
			if status.Exit != 0 {
				return fmt.Errorf("command '%s' failed with exit code %d", cmd, status.Exit)
			}
		}
	}

	return nil
}

func (w *Walker) getFiles(dir string) error {
	out := func(name, line string) {
		w.files <- line
	}
	err := func(name, line string) {
		w.logger.Errorw("shell", "error", line)
	}

	cmd := fmt.Sprintf("find '%s' -type f", dir)
	// w.logger.Debugf("running command: %s", cmd)
	status, e := Shell(cmd, ShellOptions{Out: out, Err: err})
	if e != nil {
		return e
	}
	if status.Exit != 0 {
		return fmt.Errorf("command '%s' failed with exit code %d", cmd, status.Exit)
	}
	return nil
}
