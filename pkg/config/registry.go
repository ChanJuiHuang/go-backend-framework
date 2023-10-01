package config

import (
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

type registry struct {
	viper   *viper.Viper
	configs map[string]any
	once    sync.Once
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

func (r *registry) Register(configs map[string]any) {
	r.once.Do(func() {
		for key, config := range configs {
			err := r.viper.UnmarshalKey(key, config)
			if err != nil {
				panic(err)
			}
			r.configs[key] = config
		}
	})
}

func (r *registry) Set(key string, config any) {
	r.configs[key] = config
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
