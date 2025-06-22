package main

import (
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectBoolean struct {
	Value bool
}

func (b *ObjectBoolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *ObjectBoolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type ObjectInteger struct {
	Value int64
}

// Inspect implements Object.
func (i *ObjectInteger) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *ObjectInteger) Type() ObjectType { return INTEGER_OBJ }

type ObjectNull struct{}

func (n *ObjectNull) Type() ObjectType { return NULL_OBJ }
func (n *ObjectNull) Inspect() string  { return "null" }
