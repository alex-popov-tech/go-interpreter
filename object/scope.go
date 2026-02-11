package object

import (
	"fmt"
	"slices"
	"strings"
)

type Scope struct {
	parent *Scope
	s      map[string]Object
}

func NewGlobalScope() *Scope {
	global := &Scope{parent: nil, s: map[string]Object{}}
	return global
}

func (me *Scope) Spawn() *Scope {
	return &Scope{parent: me, s: map[string]Object{}}
}

func (me *Scope) Get(identifier string) (Object, bool) {
	res, has := me.s[identifier]
	if has {
		return res, true
	}
	if me.parent != nil {
		return me.parent.Get(identifier)
	}
	return NULL_OBJECT, false
}

func (me *Scope) Add(identifier string, val Object) {
	me.s[identifier] = val
}

func (me *Scope) Set(identifier string, val Object) bool {
	_, has := me.s[identifier]
	if has {
		me.s[identifier] = val
		return true
	}
	if me.parent != nil {
		return me.parent.Set(identifier, val)
	}
	return false
}

func (me *Scope) String() string {
	rows := []string{}
	for pointer := me; pointer != nil; pointer = pointer.parent {
		rows = append(rows, fmt.Sprintf("%#v", pointer.s))
	}
	slices.Reverse(rows)
	for i, v := range rows {
		rows[i] = fmt.Sprintf("%s%s", strings.Repeat("  ", i), v)
	}
	return strings.Join(rows, "\n")
}
