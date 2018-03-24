package env

import (
	"fmt"
	. "types"
)

type Env struct {
	outer EnvType
	data  map[string]MalType
}

func NewEnv() EnvType {
	return &Env{data: make(map[string]MalType)}
}

func (env *Env) New(binds, exprs []MalType) (EnvType, error) {
	inner := Env{outer: env, data: make(map[string]MalType)}
	for i := 0; i < len(binds); i++ {
		sym, err := GetSymbol(binds[i])
		if err != nil {
			return nil, err
		}
		if sym.Value == "&" {
			sym, err := GetSymbol(binds[i+1])
			if err != nil {
				return nil, err
			}
			inner.Set(sym.Value, NewList(exprs[i:]))
			break
		}
		inner.Set(sym.Value, exprs[i])
	}
	return &inner, nil
}

func (env *Env) Set(key string, val MalType) {
	if env.data == nil {
		env.data = make(map[string]MalType)
	}
	env.data[key] = val
}

func (env *Env) Find(key string) EnvType {
	if _, ok := env.data[key]; ok {
		return env
	} else if env.outer != nil {
		return env.outer.Find(key)
	} else {
		return nil
	}
}

func (env *Env) Get(key string) (MalType, error) {
	e := env.Find(key)
	if e == nil {
		return nil, fmt.Errorf("unknown key: %v", key)
	}
	return e.(*Env).data[key], nil
}
