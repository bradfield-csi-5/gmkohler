package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

// Given an AST node corresponding to a function (guaranteed to be
// of the form `func f(x, y byte) byte`), compile it into assembly
// code.
//
// Recall from the README that the input parameters `x` and `y` should
// be read from memory addresses `1` and `2`, and the return value
// should be written to memory address `0`.
func compile(node *ast.FuncDecl) (string, error) {
	var buf strings.Builder
	paramNames := node.Type.Params.List[0].Names
	params := make([]string, len(paramNames))
	for j, p := range paramNames {
		params[j] = p.Name
	}
	for _, s := range node.Body.List {
		switch st := s.(type) {
		case *ast.ReturnStmt:
			res := st.Results[0]
			switch rt := res.(type) {
			case *ast.BasicLit:
				intVal, err := strconv.Atoi(rt.Value)
				if err != nil {
					return "", err
				}
				retVal := byte(intVal)
				buf.WriteString(fmt.Sprintf("pushi %d\n", retVal))
			case *ast.BinaryExpr:
				xStatement, err := handleParam(rt.X, params)
				if err != nil {
					return "", err
				}
				buf.WriteString(xStatement)
				yStatement, err := handleParam(rt.Y, params)
				if err != nil {
					return "", err
				}
				buf.WriteString(yStatement)
				switch rt.Op {
				case token.ADD:
					buf.WriteString("add\n")
				case token.SUB:
					buf.WriteString("sub\n")
				case token.MUL:
					buf.WriteString("mul\n")
				case token.QUO:
					buf.WriteString("div\n")
				default:
					return "", fmt.Errorf("unsupported token %s", rt.Op.String())
				}
			default:
				return "", fmt.Errorf(
					"unsupported return type %T (%+v)",
					rt,
					rt,
				)
			}
			buf.WriteString("pop 0\n")
		default:
			return "", fmt.Errorf("unsupported statement type %T (%+v)", st, st)
		}
	}
	fmt.Printf("%T", node)
	buf.WriteString("halt\n")
	return buf.String(), nil
}
func handleParam(ident ast.Expr, params []string) (string, error) {
	switch xt := ident.(type) {
	case *ast.BasicLit:
		fmt.Printf("found type %T", xt)
		v, err := strconv.Atoi(xt.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("pushi %d\n", v), nil
	case *ast.Ident:
		j := slices.Index(params, xt.Name)
		if j < 0 {
			return "", fmt.Errorf(
				"found identifier %s not in parameters", xt.Name)
		}
		// params are 1-indexed
		return fmt.Sprintf("push %d\n", j+1), nil
	default:
		fmt.Printf("found type %T", xt)
		return "", nil
	}
}
