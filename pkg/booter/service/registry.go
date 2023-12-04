package service

import "reflect"

type registry struct {
	services map[string]any
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

func (r *registry) Set(key string, service any) {
	if !(reflect.ValueOf(service).Kind() == reflect.Pointer) {
		panic("service is not the pointer")
	}
	r.services[key] = service
}

func (r *registry) SetMany(services map[string]any) {
	for key, service := range services {
		r.Set(key, service)
	}
}

func (r *registry) Get(key string) any {
	v := reflect.ValueOf(r.services[key])

	return v.Interface()
}

func (r *registry) Clone(key string) any {
	v := reflect.ValueOf(r.services[key])

	return v.Elem().Interface()
}
