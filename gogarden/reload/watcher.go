package reload

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/maps"
)

type watcher struct {
	watcher *fsnotify.Watcher

	wakeup chan struct{}

	mu      sync.Mutex
	watches map[string]bool
	pending map[string]bool
}

func newWatcher() *watcher {
	return &watcher{
		wakeup:  make(chan struct{}, 1),
		watches: map[string]bool{},
		pending: map[string]bool{},
	}
}

func (r *watcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	r.watcher = watcher

	go r.process()

	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			r.addWatch(path)
		}
		return nil
	})
	if err != nil {
		watcher.Close()
		return err
	}

	return nil
}

func (r *watcher) Wait(ctx context.Context, debounce time.Duration) error {
	select {
	case <-r.wakeup:
	case <-ctx.Done():
		return ctx.Err()
	}

	time.Sleep(debounce)

	select {
	case <-r.wakeup:
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}

func (r *watcher) Changes() []string {
	select {
	case <-r.wakeup:
	default:
	}

	r.mu.Lock()
	ret := r.pending
	r.pending = map[string]bool{}
	defer r.mu.Unlock()

	return maps.Keys(ret)
}

func (r *watcher) process() {
	for {
		select {
		case event, ok := <-r.watcher.Events:
			if !ok {
				return
			}
			event.Name = filepath.Clean(event.Name)
			if err := r.event(event); err != nil {
				log.Print(err)
			}
		case err, ok := <-r.watcher.Errors:
			if !ok {
				return
			}
			log.Print(err)
		}
	}
}

func (r *watcher) addWatch(path string) {
	path = filepath.Clean(path)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.watches[path] = true
	r.watcher.Add(path)
}

func (r *watcher) rmWatch(path string) {
	path = filepath.Clean(path)

	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.watches, path)
	r.watcher.Remove(path)
}

func (r *watcher) noteChange(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pending[path] = true
	select {
	case r.wakeup <- struct{}{}:
	default:
	}
}

func (r *watcher) event(e fsnotify.Event) error {
	if e.Has(fsnotify.Create) {
		st, err := os.Stat(e.Name)
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		} else if err != nil {
			return err
		}
		if st.IsDir() {
			r.addWatch(e.Name)
		} else {
			r.noteChange(e.Name)
		}
	}
	if e.Has(fsnotify.Write) {
		r.noteChange(e.Name)
	}
	if e.Has(fsnotify.Remove) || e.Has(fsnotify.Rename) {
		if r.watches[e.Name] {
			r.rmWatch(e.Name)
		} else {
			r.noteChange(e.Name)
		}
	}
	return nil
}
