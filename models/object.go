package models

import (
	"bytes"
	"fmt"
	"github.com/kanersps/loop/ast"
	"hash/fnv"
	"strings"
)

const (
	INTEGER   = "INTEGER"
	BOOLEAN   = "BOOLEAN"
	TYPE_NULL = "NULL"
	RETURN    = "RETURN"
	ERROR     = "ERROR"
	FUNCTION  = "FUNCTION"
	STRING    = "STRING"
	BUILTIN   = "BUILTIN"
	ARRAY     = "ARRAY"
	HASH      = "HASH"
)

type ObjectType string

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

type Hashable interface {
	// TODO: Cache return values of the HashKey method to improve performance
	HashKey() HashKey
}

func (h *Hash) Type() ObjectType {
	return HASH
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	} else {
		return HashKey{Type: b.Type(), Value: 0}
	}
}

type Null struct{}

func (n *Null) Type() ObjectType { return TYPE_NULL }
func (n *Null) Inspect() string  { return "null" }

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURN }
func (r *Return) Inspect() string  { return r.Inspect() }

// TODO: Implement line number & column number to more easily debug issues
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "Exception: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string {
	return s.Value
}
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Func BuiltinFunction
	Env  *Environment
}

func (b *Builtin) Type() ObjectType { return BUILTIN }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (array *Array) Type() ObjectType { return ARRAY }
func (array *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range array.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
