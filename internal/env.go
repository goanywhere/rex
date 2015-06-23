package internal

import (
	"path"

	"github.com/goanywhere/env"
)

type Env struct {
	Base string
}

func New(dotenv string) *Env {
	env.Load(dotenv)
	return &Env{Base: path.Base(dotenv)}
}

func (self *Env) Get(key string) (string, bool) {
	return env.Get(key)
}

func (self *Env) Set(key string, value interface{}) error {
	return env.Set(key, value)
}

func (self *Env) String(key string, fallback ...string) string {
	return env.String(key, fallback...)
}

func (self *Env) Strings(key string, fallback ...[]string) []string {
	return env.Strings(key, fallback...)
}

func (self *Env) Int(key string, fallback ...int) int {
	return env.Int(key, fallback...)
}

func (self *Env) Int64(key string, fallback ...int64) int64 {
	return env.Int64(key, fallback...)
}

func (self *Env) Uint(key string, fallback ...uint) uint {
	return env.Uint(key, fallback...)
}

func (self *Env) Uint64(key string, fallback ...uint64) uint64 {
	return env.Uint64(key, fallback...)
}

func (self *Env) Bool(key string, fallback ...bool) bool {
	return env.Bool(key, fallback...)
}

func (self *Env) Float(key string, fallback ...float64) float64 {
	return env.Float(key, fallback...)
}

func (self *Env) Map(spec interface{}) error {
	return env.Map(spec)
}
