package fsmonitor

import (
	"code.google.com/p/go.exp/fsnotify"
  "github.com/datastream/btree"
	"os"
	"path/filepath"
)

func NewWatcher() (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	monitorWatcher := initWatcher(watcher, []string{})
	return monitorWatcher, nil
}

func NewWatcherWithSkipFolders(skipFolders []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	monitorWatcher := initWatcher(watcher, skipFolders)
	return monitorWatcher, nil
}

func initWatcher(watcher *fsnotify.Watcher, skipFolders []string) *Watcher {
	event := make(chan *fsnotify.FileEvent)
	watcherError := make(chan error)
  tree := btree.NewBtree()
	monitorWatcher := &Watcher{Event: event, Error: watcherError, watcher: watcher, SkipFolders: skipFolders, watchTree: tree}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				event <- ev
				if ev.IsCreate() {
					go func() {
						if f, err := os.Stat(ev.Name); err == nil {
							if f.IsDir() {
								monitorWatcher.watchAllFolders(ev.Name)
							}
						}

					}()
				}
				if ev.IsDelete() {
					go func() {
            monitorWatcher.Remove(ev.Name)
						// watcher.RemoveWatch(ev.Name)
					}()
				}
			case e := <-watcher.Error:
				watcherError <- e
			}
		}
	}()
	return monitorWatcher
}

type Watcher struct {
	Event       chan *fsnotify.FileEvent
	Error       chan error
	SkipFolders []string
	watcher     *fsnotify.Watcher
  watchDir    *btree.Btree
}

func (self *Watcher) Watch(path string) error {
	err := self.watchAllFolders(path)
	if err != nil {
		return err
	}
	return nil
}

func (self *Watcher) watchAllFolders(path string) (err error) {
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f != nil && f.IsDir() {
			filename := f.Name()
			for _, skipFolder := range self.SkipFolders {
				match, err := filepath.Match(skipFolder, filename)
				if err != nil {
					return err
				}
				if match {
					return filepath.SkipDir
				}
			}
      err := self.watchTree.Insert([]byte(path), []byte(path))
      if err != nil {
        return err
      }
			err := self.addWatcher(path)
			if err != nil {
        // 失敗したらTreeから取り除く
        self.watchTree.Delete([]byte(path))
				return err
			}
		}
		return nil
	})
	return
}

func (self *Watcher) addWatcher(path string) (err error) {
	err = self.watcher.Watch(path)
	return
}

func (self *Watcher) Remove(path string) (err error) {
  err = self.watcher.RemoveWatch(path)
  return
}
