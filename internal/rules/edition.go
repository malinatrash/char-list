package rules

import (
	"encoding/json"
	"io"
	"slices"

	"char-list/internal/dice"
)

type Edition interface {
	ID() string
	Name() string
	Blurb() string
	Method() string
	Data() any
	Random(raw json.RawMessage, fields []string, method string) (any, []dice.Log, error)
	Derive(raw json.RawMessage) (any, error)
	PDF(raw json.RawMessage, w io.Writer) error
}

type Registry struct {
	list []Edition
	byID map[string]Edition
}

func NewRegistry(eds ...Edition) *Registry {
	r := &Registry{byID: map[string]Edition{}}
	for _, e := range eds {
		r.list = append(r.list, e)
		r.byID[e.ID()] = e
	}
	return r
}

func (r *Registry) All() []Edition { return r.list }

func (r *Registry) Get(id string) (Edition, bool) {
	e, ok := r.byID[id]
	return e, ok
}

func has(fields []string, f string) bool {
	if fields == nil {
		return true
	}
	return slices.Contains(fields, f)
}

func contains(xs []string, s string) bool {
	return slices.Contains(xs, s)
}
