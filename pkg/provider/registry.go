package provider

import (
	"reflect"
	"sync"
)

type registry struct {
	services map[string]any
	once     sync.Once
}

var Registry *registry

func init() {
	Registry = &registry{
		services: map[string]any{},
	}
}

func NewRegistry() *registry {
	return &registry{
		services: map[string]any{},
	}
}

func (r *registry) Register(services map[string]any) {
	r.once.Do(func() {
		for key, service := range services {
			r.services[key] = service
		}
	})
}

func (r *registry) Set(key string, service any) {
	r.services[key] = service
}

func (r *registry) Get(key string) any {
	v := reflect.ValueOf(r.services[key])

	return v.Interface()
}

func (r *registry) Clone(key string) any {
	v := reflect.ValueOf(r.services[key])

	return v.Elem().Interface()
}
