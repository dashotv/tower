package app

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func (s *Server) Watcher() {
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		s.Log.Fatal("setting up filesystem watcher: ", err)
	}

	size := 25
	ch := make(chan string)
	for i := 0; i < size; i++ {
		go func(j int) {
			for path := range ch {
				s.Log.Infof("filesystem watcher:%2d %s", j, path)
				err := s.watcher.Add(path)
				if err != nil {
					s.Log.Warnf("error adding %s: %s", path, err)
				}
			}
		}(i)
	}

	total := 0
	for _, fs := range ConfigInstance().Filesystems.Directories {
		_ = filepath.Walk(fs, func(path string, fi os.FileInfo, err error) error {
			if fi.IsDir() {
				ch <- path
				total += 1
			}
			return nil
		})
	}
	s.Log.Infof("filesystem watcher: watching %d directories", total)

	go func() {
		s.Log.Info("starting filesystem watcher...")
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					s.Log.Printf("filesystem watcher: event not ok")
					return
				}
				s.Log.Printf("%s %s", event.Name, event.Op)
			case err, ok := <-s.watcher.Errors:
				if !ok {
					s.Log.Printf("filesystem watcher: error not ok")
					return
				}
				s.Log.Printf("error watching files: %s", err)
			}
		}
	}()
}
