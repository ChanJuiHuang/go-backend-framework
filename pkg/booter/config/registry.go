package config

import (
	"reflect"

	"github.com/spf13/viper"
)

type registry struct {
	viper   *viper.Viper
	configs map[string]any
}

var Registry *registry

func init() {
	Registry = &registry{
		configs: map[string]any{},
	}
}

func NewRegistry(v *viper.Viper) *registry {
	return &registry{
		viper:   v,
		configs: map[string]any{},
	}
}

func (r *registry) Set(key string, config any) {
	if !(reflect.ValueOf(config).Kind() == reflect.Pointer) {
		panic("config is not the pointer")
	}
	r.configs[key] = config
}

func (r *registry) SetMany(configs map[string]any) {
	for key, config := range configs {
		r.Set(key, config)
	}
}

func (r *registry) Register(key string, config any) {
	err := r.viper.UnmarshalKey(key, config)
	if err != nil {
		panic(err)
	}
	r.Set(key, config)
}

func (r *registry) RegisterMany(configs map[string]any) {
	for key, config := range configs {
		r.Register(key, config)
	}
}

func (r *registry) Get(key string) any {
	v := reflect.ValueOf(r.configs[key])

	return v.Elem().Interface()
}

func (r *registry) SetViper(v *viper.Viper) {
	r.viper = v
}

func (r *registry) GetViper() viper.Viper {
	return *r.viper
}
