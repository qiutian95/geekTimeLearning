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
	defer providersMu.Unlock() // 和下面的直接unlock不同的是：1.执行时机不同：unlock是立即释放锁 2.用途不同：defer是在函数结束时清理工作，而unlock是显示释放锁
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
