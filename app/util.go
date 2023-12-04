package app

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var pathQuoteRegex = regexp.MustCompile(`'(\w{1,2})`)
var pathCharRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func path(title string) string {
	var s string
	s = pathQuoteRegex.ReplaceAllString(title, "$1")
	s = strings.ToLower(s)
	s = pathCharRegex.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func sumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "failed to hash file")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func sumFiles(source, destination string) (bool, error) {
	type response struct {
		Hash string
		Err  error
	}
	ch_source := make(chan response)
	ch_dest := make(chan response)

	go func() {
		sum, err := sumFile(source)
		ch_source <- response{sum, err}
	}()
	go func() {
		sum, err := sumFile(destination)
		ch_dest <- response{sum, err}
	}()

	r_source := <-ch_source
	r_dest := <-ch_dest

	if r_source.Err != nil {
		return false, errors.Wrap(r_source.Err, "failed to sum source")
	}
	if r_dest.Err != nil {
		return false, errors.Wrap(r_dest.Err, "failed to sum destination")
	}

	return r_source.Hash == r_dest.Hash, nil
}

// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file
// Copy copies the contents of the file at srcpath to a regular file at dstpath.
// If dstpath already exists and is not a directory, the function truncates it.
// The function does not copy file modes or file attributes.
func FileCopy(srcpath, dstpath string) (err error) {
	r, err := os.Open(srcpath)
	if err != nil {
		return errors.Wrap(err, "source")
	}
	defer r.Close() // ignore error: file was opened read-only.

	w, err := os.Create(dstpath)
	if err != nil {
		return errors.Wrap(err, "destination")
	}

	defer func() {
		// Report the error from Close, if any,
		// but do so only if there isn't already
		// an outgoing error.
		if c := w.Close(); c != nil && err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// WithTimeout runs a delegate function with a timeout,
//
// Example: Wait for a channel
//
//	if value, ok := WithTimeout(func()interface{}{return <- inbox}, time.Second); ok {
//	    // returned
//	} else {
//	    // didn't return
//	}
//
// Example: To send to a channel
//
//	_, ok := WithTimeout(func()interface{}{outbox <- myValue; return nil}, time.Second)
//	if !ok {
//	    // didn't send
//	}
func WithTimeout(delegate func() interface{}, timeout time.Duration) (ret interface{}, ok bool) {
	ch := make(chan interface{}, 1) // buffered
	go func() { ch <- delegate() }()
	select {
	case ret = <-ch:
		return ret, true
	case <-time.After(timeout):
	}
	return nil, false
}

func shouldDownloadFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	ext = ext[1:]
	list := lo.Filter(cfg.Extensions(), func(s string, i int) bool {
		return s == ext
	})
	return len(list) > 0
}
