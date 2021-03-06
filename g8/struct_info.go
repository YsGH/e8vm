package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/toposort"
)

type structInfo struct {
	name *lex8.Token
	ast  *ast.Struct
	t    *types.Struct // the struct type
	pt   *types.Pointer

	methodObjs []*objFunc

	deps []string
}

func addStructDeps(deps []string, t ast.Expr) []string {
	switch t := t.(type) {
	case *ast.Operand:
		if t.Token.Type == parse.Ident {
			deps = append(deps, t.Token.Lit)
		}
	case *ast.ParenExpr:
		deps = addStructDeps(deps, t.Expr)
	case *ast.ArrayTypeExpr:
		if t.Len != nil {
			// not a slice
			deps = addStructDeps(deps, t.Type)
		}
	}

	return deps
}

func structDeps(s *ast.Struct) []string {
	var ret []string
	for _, f := range s.Fields {
		ret = addStructDeps(ret, f.Type)
	}
	return ret
}

func (info *structInfo) Name() string {
	return info.name.Lit
}

func newStructInfo(s *ast.Struct) *structInfo {
	ret := new(structInfo)
	ret.name = s.Name
	ret.ast = s
	ret.deps = structDeps(s)
	ret.t = types.NewStruct(ret.name.Lit)
	ret.pt = types.NewPointer(ret.t)

	return ret
}

func sortStructs(b *builder, m map[string]*structInfo) []*structInfo {
	s := toposort.NewSorter("struct")

	for name, info := range m {
		s.AddNode(name, info.name, info.deps)
	}

	order := s.Sort(b)
	var ret []*structInfo
	for _, name := range order {
		ret = append(ret, m[name])
	}

	return ret
}
