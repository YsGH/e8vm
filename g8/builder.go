package g8

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// builder builds a package
type builder struct {
	*lex8.ErrorList
	path string

	p         *ir.Pkg
	f         *ir.Func
	fretNamed bool
	fretRef   *ref

	golike bool

	b            *ir.Block
	scope        *sym8.Scope
	symPkg       *sym8.Pkg
	structFields map[*types.Struct]*sym8.Table

	continues *blockStack
	breaks    *blockStack

	exprFunc  func(b *builder, expr ast.Expr) *ref
	stmtFunc  func(b *builder, stmt ast.Stmt)
	typeFunc  func(b *builder, expr ast.Expr) types.T
	constFunc func(b *builder, expr ast.Expr) *ref
	irLog     io.WriteCloser

	panicFunc ir.Ref

	// this pointer, only valid when building a method.
	this *ref

	// file level dependency, for checking circular dependencies.
	deps deps

	anonyCount int

	rand *rand.Rand
}

func newRand() *rand.Rand {
	var buf [8]byte
	_, err := crand.Read(buf[:])
	if err != nil {
		panic(err)
	}
	seed := int64(binary.LittleEndian.Uint64(buf[:]))
	return rand.New(rand.NewSource(seed))
}

func newBuilder(path string, golike bool) *builder {
	ret := new(builder)
	ret.ErrorList = lex8.NewErrorList()
	ret.path = path
	ret.p = ir.NewPkg(path)
	ret.scope = sym8.NewScope() // package scope
	ret.symPkg = &sym8.Pkg{path}
	ret.golike = golike

	ret.continues = newBlockStack()
	ret.breaks = newBlockStack()
	ret.structFields = make(map[*types.Struct]*sym8.Table)

	ret.rand = newRand()

	return ret
}

func (b *builder) refSym(sym *sym8.Symbol, pos *lex8.Pos) {
	// track file dependencies inside a package
	if b.deps == nil {
		return // no need to track deps
	}

	symPos := sym.Pos
	if symPos == nil {
		return // builtin
	}
	if sym.Pkg() != b.symPkg {
		return // cross package reference
	}
	if pos.File == symPos.File {
		return
	}

	b.deps.add(pos.File, symPos.File)
}

func (b *builder) anonyName(name string) string {
	if name == "_" {
		name = fmt.Sprintf("_:%d", b.anonyCount)
		b.anonyCount++
	}
	return name
}

func (b *builder) newTempIR(t types.T) ir.Ref {
	return b.f.NewTemp(t.Size(), types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) newTemp(t types.T) *ref { return newRef(t, b.newTempIR(t)) }

func (b *builder) newCond() ir.Ref { return b.f.NewTemp(1, true, false) }
func (b *builder) newPtr() ir.Ref  { return b.f.NewTemp(4, true, true) }

func (b *builder) newAddressableTemp(t types.T) *ref {
	return newAddressableRef(t, b.newTempIR(t))
}

func (b *builder) newLocal(t types.T, name string) ir.Ref {
	return b.f.NewLocal(t.Size(), name,
		types.IsByte(t), t.RegSizeAlign(),
	)
}

func (b *builder) newGlobalVar(t types.T, name string) ir.Ref {
	name = b.anonyName(name)
	return b.p.NewGlobalVar(t.Size(), name, types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) buildExpr(expr ast.Expr) *ref {
	if b.exprFunc != nil {
		return b.exprFunc(b, expr)
	}
	return nil
}

func (b *builder) buildConstExpr(expr ast.Expr) *ref {
	if b.constFunc != nil {
		return b.constFunc(b, expr)
	}
	return nil
}

func (b *builder) buildType(expr ast.Expr) types.T {
	if b.typeFunc != nil {
		return b.typeFunc(b, expr)
	}
	return nil
}

func (b *builder) buildStmts(stmts []ast.Stmt) {
	if b.stmtFunc == nil {
		return
	}

	for _, stmt := range stmts {
		b.stmtFunc(b, stmt)
	}
}

func (b *builder) buildStmt(stmt ast.Stmt) { b.stmtFunc(b, stmt) }
