package sym8

import (
	"e8vm.io/e8vm/lex8"
)

// Symbol is a data structure for saving a symbol.
type Symbol struct {
	pkg  *Pkg
	name string

	Type int
	Item interface{}
	Pos  *lex8.Pos
}

// Name returns the symbol name.
// This name is immutable for its used for indexing in the tables.
func (s *Symbol) Name() string { return s.name }

// Pkg returns the package token of the symbol.
func (s *Symbol) Pkg() *Pkg { return s.pkg }

// Clone returns a clone of the symbol which item is the given one.
func (s *Symbol) Clone(item interface{}) *Symbol {
	ret := *s
	ret.Item = item
	return &ret
}

// Make creates a new symbol
func Make(
	pkg *Pkg,
	name string,
	t int,
	item interface{},
	pos *lex8.Pos,
) *Symbol {
	return &Symbol{pkg, name, t, item, pos}
}
