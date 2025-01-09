package runner

import (
	"github.com/fsnotify/fsnotify"
	core "github.com/xedeveloper/Paradigm/engine/core/component"
)

type HotReload struct {
	watcher    *fsnotify.Watcher
	components map[string]*core.Component
}

func NewHotReload() (*HotReload, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &HotReload{
		watcher:    watcher,
		components: make(map[string]*core.Component),
	}, nil
}

func (h *HotReload) Watch(paths ...string) error {
	for _, path := range paths {
		if err := h.watcher.Add(path); err != nil {
			return err
		}
	}
	return nil
}
