package main

import (
	"errors"
	"path"
	"github.com/fsnotify/fsnotify"
)

type WatcherHandler interface {
	CanHandle(path string) bool
	Handle(ev *fsnotify.Event)
}

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	handlers []WatcherHandler
	started  bool
}

func (fw *FileWatcher) AddWatch(path string) error {
	return fw.watcher.Add(path)
}

func (fw *FileWatcher) RemoveWatch(strFile string) error {
	return fw.watcher.Remove(strFile)
}

func (fw *FileWatcher) AddHandler(detector WatcherHandler) error {
	if detector == nil {
		return errors.New("The parameter 'detector' cannot be nil")
	}
	fw.handlers = append(fw.handlers, detector)
	return nil
}

func (fw *FileWatcher) Start() {
	if fw.started {
		return
	}
	fw.started = true
	go func() {
		for {
			select {
			case ev := <-fw.watcher.Events:
				for _, detector := range fw.handlers {
					if detector.CanHandle(path.Clean(ev.Name)) {
						detector.Handle(&ev)
					}
				}
			}
		}
	}()
}

func NewWatcher() (*FileWatcher, error) {
	tmpWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &FileWatcher{
		watcher:  tmpWatcher,
	}
	return w, nil
}