package asm8

import (
	"e8vm.io/e8vm/link8"
	"e8vm.io/e8vm/sym8"
)

// Lib is the compiler output of a package
// it contains the package for linking,
// and also the symbols for importing
type lib struct {
	*link8.Pkg

	symbols map[string]*sym8.Symbol
}

func (p *lib) Main() string { return "main" }

func (p *lib) Tests() (map[string]uint32, string) {
	// Assembly now does not have tests.
	return nil, ""
}

// NewPkgObj creates a new package compile object
func newLib(p string) *lib {
	ret := new(lib)
	ret.Pkg = link8.NewPkg(p)
	ret.symbols = make(map[string]*sym8.Symbol)
	return ret
}

// Link returns the link8.Package for linking.
func (p *lib) Link() *link8.Pkg { return p.Pkg }

func (p *lib) Declare(s *sym8.Symbol) {
	_, found := p.symbols[s.Name()]
	if found {
		panic("redeclare")
	}
	p.symbols[s.Name()] = s

	switch s.Type {
	case SymConst:
		panic("todo")
	case SymFunc:
		p.Pkg.DeclareFunc(s.Name())
	case SymVar:
		p.Pkg.DeclareVar(s.Name())
	default:
		panic("declare with invalid sym type")
	}
}

// Query returns the symbol declared by name and its symbol index
// if the symbol is a function or variable. It returns nil, 0 when
// the symbol of name is not found.
func (p *lib) query(name string) *sym8.Symbol {
	ret, found := p.symbols[name]
	if !found {
		return nil
	}

	switch ret.Type {
	case SymConst:
		return ret
	case SymFunc, SymVar:
		s := p.Pkg.SymbolByName(name)
		if s == nil {
			panic("symbol missing")
		}
		return ret
	default:
		panic("bug")
	}
}

// Lib retunrs the linkable lib.
func (p *lib) Lib() *link8.Pkg { return p.Pkg }

// Symbols returns "asm8", nil. Linking with assembly should directly look
// into the lib.
func (p *lib) Symbols() (string, *sym8.Table) { return "asm8", nil }
