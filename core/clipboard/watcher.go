package clipboard

import (
	"fmt"

	"github.com/google/uuid"
)

type Watcher struct {
	listeners map[string]func(data []byte)
}

func NewWatcher() *Watcher {
	return &Watcher{
		listeners: make(map[string]func(data []byte)),
	}
}

func (w *Watcher) Add_listener(listener func(data []byte)) (string, func()) {
	key := uuid.New().String()
	fmt.Println("Adding listener with key", key)
	w.listeners[key] = listener

	return key, func() {
		delete(w.listeners, key)
	}
}

func (w *Watcher) Remove_listener(key string) {
	fmt.Println("Removing listener with key", key)
	delete(w.listeners, key)
}

func (w *Watcher) Notify(data []byte) {
	for _, listener := range w.listeners {
		listener(data)
	}
}

func (w *Watcher) Remove_all_listeneres() {
	w.listeners = make(map[string]func(data []byte))
}
