package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
)

// Given an expression containing only int types, evaluate
// the expression and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		x, err := Evaluate(e.X)
		if err != nil {
			return 0, err
		}
		y, err := Evaluate(e.Y)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			return x / y, nil
		}
	case *ast.BasicLit:
		n, err := strconv.Atoi(e.Value)
		if err != nil {
			return 0, err
		}
		return n, nil
	case *ast.ParenExpr:
		return Evaluate(e.X)
	default:
		return 0, fmt.Errorf("unsupported expression type %T", e)
	}
	return 0, nil
}

func main() {
	var expr ast.Expr
	expr, err := parser.ParseExpr("1 + 2 - 3 * 4")
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	err = ast.Print(fset, expr)
	if err != nil {
		log.Fatal(err)
	}
}
