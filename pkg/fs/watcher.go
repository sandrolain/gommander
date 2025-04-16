package fs

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type DirWatcherCallback func(string, error)

type DirWatcher struct {
	watcher       *fsnotify.Watcher
	stop          chan struct{}
	subscriptions []*DirWatcherSubscription
	sLock         sync.Mutex
}

func NewDirWatcher() (*DirWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &DirWatcher{
		watcher:       watcher,
		stop:          make(chan struct{}),
		subscriptions: make([]*DirWatcherSubscription, 0),
	}, nil
}

func (dw *DirWatcher) Watch(path string) error {
	go func() {
		defer dw.watcher.Close()
		for {
			select {
			case event, ok := <-dw.watcher.Events:
				if !ok {
					return
				}
				dw.notifySubscribers(event.Name, nil)
			case err, ok := <-dw.watcher.Errors:
				if !ok {
					return
				}
				dw.notifySubscribers("", err)
			case <-dw.stop:
				return
			}
		}
	}()

	err := dw.watcher.Add(path)
	if err != nil {
		close(dw.stop)
		return err
	}

	return nil
}

func (dw *DirWatcher) Subscribe(path string, callback DirWatcherCallback) (*DirWatcherSubscription, error) {
	dw.sLock.Lock()
	defer dw.sLock.Unlock()

	s := &DirWatcherSubscription{
		ID:       time.Now().UnixMilli(),
		Path:     path,
		Callback: callback,
	}

	dw.subscriptions = append(dw.subscriptions, s)
	return s, nil
}

func (dw *DirWatcher) Unsubscribe(s *DirWatcherSubscription) bool {
	dw.sLock.Lock()
	defer dw.sLock.Unlock()

	for i, sub := range dw.subscriptions {
		if sub.ID == s.ID {
			dw.subscriptions = append(dw.subscriptions[:i], dw.subscriptions[i+1:]...)
			break
		}
	}

	return len(dw.subscriptions) == 0
}

func (dw *DirWatcher) notifySubscribers(eventPath string, err error) {
	dw.sLock.Lock()
	defer dw.sLock.Unlock()

	for _, s := range dw.subscriptions {
		s.Callback(eventPath, err)
	}
}

func (dw *DirWatcher) Stop() {
	dw.sLock.Lock()
	defer dw.sLock.Unlock()

	close(dw.stop)
	dw.subscriptions = []*DirWatcherSubscription{}
}

type DirWatcherSubscription struct {
	ID       int64
	Path     string
	Callback DirWatcherCallback
}

var watchers map[string]*DirWatcher
var watchersLock sync.Mutex

func SubscribeWatcher(prevSub *DirWatcherSubscription, path string, callback DirWatcherCallback) (*DirWatcherSubscription, error) {
	watchersLock.Lock()
	defer watchersLock.Unlock()

	if watchers == nil {
		watchers = make(map[string]*DirWatcher)
	}

	if prevSub != nil {
		if watcher, exists := watchers[prevSub.Path]; exists {
			if watcher.Unsubscribe(prevSub) && path != prevSub.Path {
				watcher.Stop()
				delete(watchers, prevSub.Path)
			}
		}
	}

	// Check if a watcher already exists for the path
	if watcher, exists := watchers[path]; exists {
		return watcher.Subscribe(path, callback)
	}

	// Create a new watcher if it doesn't exist
	newWatcher, err := NewDirWatcher()
	if err != nil {
		return nil, err
	}

	err = newWatcher.Watch(path)
	if err != nil {
		return nil, err
	}

	// Store the new watcher in the map
	watchers[path] = newWatcher

	// Subscribe to the new watcher
	return newWatcher.Subscribe(path, callback)
}
