package factory

import (
	"bookstore/store"
	"fmt"
	"sync"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]store.Store)
)

func Register(name string, p store.Store) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if p == nil {
		panic("store：Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("store: Register called twice for provider " + name)
	}
	providers[name] = p
}

func New(providerName string) (store.Store, error) {
	providersMu.Lock()
	p, ok := providers[providerName]
	providersMu.Unlock()
	if !ok {
		return nil, fmt.Errorf("store: provider %s not exists", providerName)
	}
	return p, nil

}
