package main

import "sync"

type handleRegistry struct {
	mu     sync.Mutex
	next   uintptr
	values map[uintptr]any
}

func newHandleRegistry() *handleRegistry {
	return &handleRegistry{
		next:   1,
		values: make(map[uintptr]any),
	}
}

func (r *handleRegistry) register(value any) uintptr {
	r.mu.Lock()
	defer r.mu.Unlock()

	handle := r.next
	r.next++
	r.values[handle] = value
	return handle
}

func (r *handleRegistry) take(handle uintptr) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	value, ok := r.values[handle]
	if !ok {
		return nil, false
	}
	delete(r.values, handle)
	return value, true
}

func (r *handleRegistry) get(handle uintptr) (any, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	value, ok := r.values[handle]
	return value, ok
}

