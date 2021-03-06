package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildStmt(b *builder, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.EmptyStmt:
		// do nothing
	case *ast.ExprStmt:
		buildExprStmt(b, stmt.Expr)
	case *ast.IncStmt:
		buildIncStmt(b, stmt)
	case *ast.DefineStmt:
		buildDefineStmt(b, stmt)
	case *ast.AssignStmt:
		buildAssignStmt(b, stmt)
	case *ast.IfStmt:
		buildIfStmt(b, stmt)
	case *ast.ForStmt:
		buildForStmt(b, stmt)
	case *ast.BlockStmt:
		buildBlock(b, stmt.Block)
	case *ast.VarDecls:
		buildVarDecls(b, stmt)
	case *ast.ConstDecls:
		buildConstDecls(b, stmt)
	case *ast.ReturnStmt:
		buildReturnStmt(b, stmt)
	case *ast.ContinueStmt:
		buildContinueStmt(b, stmt)
	case *ast.BreakStmt:
		buildBreakStmt(b, stmt)
	default:
		b.Errorf(nil, "invalid or not implemented: %T", stmt)
	}
}
