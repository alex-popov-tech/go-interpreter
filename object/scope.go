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

func (self *Scope) Spawn() *Scope {
	return &Scope{parent: self, s: map[string]Object{}}
}

func (self *Scope) Get(identifier string) (Object, bool) {
	for pointer := self; pointer != nil; pointer = pointer.parent {
		res, isOk := pointer.s[identifier]
		if isOk {
			return res, true
		}
	}
	return NULL_OBJECT, false
}

func (self *Scope) Add(identifier string, val Object) {
	self.s[identifier] = val
}

func (self *Scope) Set(identifier string, val Object) bool {
	for pointer := self; pointer != nil; pointer = pointer.parent {
		_, isExist := pointer.s[identifier]
		if isExist {
			pointer.s[identifier] = val
			return true
		}
	}
	return false
}

func (self *Scope) String() string {
	rows := []string{}
	for pointer := self; pointer != nil; pointer = pointer.parent {
		rows = append(rows, fmt.Sprintf("%#v", pointer.s))
	}
	slices.Reverse(rows)
	for i, v := range rows {
		rows[i] = fmt.Sprintf("%s%s", strings.Repeat("  ", i), v)
	}
	return strings.Join(rows, "\n")
}
