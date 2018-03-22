package env

import (
	"fmt"
	. "types"
)

type Env struct {
	outer *Env
	data  map[string]MalType
}

func (env *Env) Set(key string, val MalType) {
	if env.data == nil {
		env.data = make(map[string]MalType)
	}
	env.data[key] = val
}

func (env *Env) Find(key string) (MalType, bool) {
	if env == nil {
		return nil, false
	}
	if val, ok := env.data[key]; ok {
		return val, true
	}
	return env.outer.Find(key)
}

func (env *Env) Get(key string) (MalType, error) {
	if val, ok := env.Find(key); ok {
		return val, nil
	}
	return nil, fmt.Errorf("unknown key: %v", key)
}

func (env *Env) Inner() Env {
	return Env{outer: env}
}
