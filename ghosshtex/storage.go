package ghosshtex

import "sync"

type SharedResource struct {
	value string
	mu    sync.RWMutex
}

func (r *SharedResource) Update(value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.value = value
}

func (r *SharedResource) Read() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.value
}

func NewSharedResource(value string) *SharedResource {
	return &SharedResource{value: value}
}
