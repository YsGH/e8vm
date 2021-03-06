package g8

import (
	"math"
	"strconv"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func buildInt(b *builder, op *lex8.Token) *ref {
	ret, e := strconv.ParseInt(op.Lit, 0, 64)
	if e != nil {
		b.Errorf(op.Pos, "invalid integer: %s", e)
		return nil
	}

	if ret < math.MinInt32 {
		b.Errorf(op.Pos, "integer too small to fit in 32-bit")
		return nil
	} else if ret > math.MaxUint32 {
		b.Errorf(op.Pos, "integer too large to fit in 32-bit")
		return nil
	}

	return newRef(types.NewNumber(ret), nil)
}

func buildChar(b *builder, op *lex8.Token) *ref {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid char: %s", e)
		return nil
	} else if len(v) != 1 {
		b.Errorf(op.Pos, "invalid char in quote: %q", v)
		return nil
	}
	return newRef(types.Int8, ir.Byt(v[0]))
}

func buildString(b *builder, op *lex8.Token) *ref {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid string: %s", e)
		return nil
	}

	ret := b.newTemp(types.String) // make a temp slice
	b.b.Arith(ret.IR(), nil, "makeStr", b.p.NewString(v))
	return ret
}

func buildField(b *builder, this ir.Ref, field *types.Field) *ref {
	retIR := ir.NewAddrRef(
		this,
		field.T.Size(),
		field.Offset(),
		types.IsByte(field.T),
		true,
	)
	return newAddressableRef(field.T, retIR)
}

func buildConstIdent(b *builder, ident *lex8.Token) *ref {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)
	switch s.Type {
	case symType, symStruct:
		return s.Item.(*objType).ref
	case symConst:
		return s.Item.(*objConst).ref
	case symImport:
		return s.Item.(*objImport).ref
	default:
		b.Errorf(ident.Pos, "%s is a %s, not a const",
			ident.Lit, symStr(s.Type),
		)
		return nil
	}
}

func buildIdent(b *builder, ident *lex8.Token) *ref {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifer %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	switch s.Type {
	case symVar:
		return s.Item.(*objVar).ref
	case symFunc:
		v := s.Item.(*objFunc)
		if !v.isMethod {
			return v.ref
		}
		if b.this == nil {
			panic("this missing")
		}
		return newRecvRef(v.Type().(*types.Func), b.this, v.IR())
	case symConst:
		return s.Item.(*objConst).ref
	case symType, symStruct:
		return s.Item.(*objType).ref
	case symField:
		v := s.Item.(*objField)
		return buildField(b, b.this.IR(), v.Field)
	case symImport:
		return s.Item.(*objImport).ref
	default:
		b.Errorf(ident.Pos, "todo: token type: %s", symStr(s.Type))
		return nil
	}
}

func buildOperand(b *builder, op *ast.Operand) *ref {
	if op.Token.Type == parse.Keyword && op.Token.Lit == "this" {
		if b.this == nil {
			b.Errorf(op.Token.Pos, "using this out of a method function")
			return nil
		}
		return b.this
	}

	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Char:
		return buildChar(b, op.Token)
	case parse.String:
		return buildString(b, op.Token)
	case parse.Ident:
		return buildIdent(b, op.Token)
	default:
		b.Errorf(op.Token.Pos, "invalid or not implemented: %d",
			op.Token.Type,
		)
		return nil
	}
}

func buildConstOperand(b *builder, op *ast.Operand) *ref {
	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Ident:
		return buildConstIdent(b, op.Token)
	default:
		b.Errorf(op.Token.Pos, "not a const")
		return nil
	}
}
