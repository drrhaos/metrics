package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  `Проверяет использование прямого вызова os.Exit в функции main пакета main`,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.FuncDecl) {
		for _, ttt := range x.Body.List {
			if exprStmt, ok := ttt.(*ast.ExprStmt); ok {
				if call, ok := exprStmt.X.(*ast.CallExpr); ok {
					if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
						funcName := fun.Sel.Name
						if iden, ok := fun.X.(*ast.Ident); ok {
							if iden.Name == "os" && funcName == "Exit" {
								pass.Reportf(iden.NamePos, "использование прямого вызова os.Exit в функции main пакета main")
							}
						}
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					expr(x)
				}
			}
			return true
		})
	}

	return nil, nil
}
