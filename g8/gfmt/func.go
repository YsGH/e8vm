package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printParaList(p *fmt8.Printer, lst *ast.ParaList) {
	fmt.Fprint(p, "(")
	for i, para := range lst.Paras {
		if i > 0 {
			fmt.Fprint(p, ", ")
		}
		if para.Ident != nil {
			fmt.Fprintf(p, "%s", para.Ident.Lit)
			if para.Type != nil {
				fmt.Fprint(p, " ")
			}
		}

		if para.Type != nil {
			printExpr(p, para.Type)
		}
	}
	fmt.Fprint(p, ")")
}

func printFunc(p *fmt8.Printer, f *ast.Func) {
	fmt.Fprintf(p, "func %s", f.Name.Lit)
	printParaList(p, f.Args)
	if f.RetType != nil {
		fmt.Fprint(p, " ")
		printExpr(p, f.RetType)
	} else if f.Rets != nil {
		fmt.Fprint(p, " ")
		printParaList(p, f.Rets)
	}

	fmt.Fprint(p, " ")
	printStmt(p, f.Body)
	fmt.Fprintln(p)
}
