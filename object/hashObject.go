package object

import (
	"fmt"
	"strings"
)

type HashObject struct {
	Map map[any]Object
}

func (ao HashObject) Inspect() string {
	args := []string{}
	for key, val := range ao.Map {
		args = append(args, fmt.Sprintf("%s:%s", key, val.Inspect()))
	}
	return fmt.Sprintf("#{ %s }", strings.Join(args, ", "))
}

func (ao HashObject) Type() ObjectType {
	return HASH
}
