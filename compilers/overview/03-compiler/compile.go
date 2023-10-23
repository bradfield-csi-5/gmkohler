package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// Given an AST node corresponding to a function (guaranteed to be
// of the form `func f(x, y byte) byte`), compile it into assembly
// code.
//
// Recall from the README that the input parameters `x` and `y` should
// be read from memory addresses `1` and `2`, and the return value
// should be written to memory address `0`.
func compile(node *ast.FuncDecl) (string, error) {
	var buf bytes.Buffer
	var nextLoc byte = 1
	params := make(map[string]byte)
	for _, p := range node.Type.Params.List {
		for _, n := range p.Names {
			params[n.Name] = nextLoc
			nextLoc++
		}
	}
	s := scope{identifiers: params}
	err := s.compileStatement(node.Body, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type scope struct {
	identifiers map[string]byte
}

func (s *scope) compileExpression(e ast.Expr, buf *bytes.Buffer) error {
	var err error
	switch et := e.(type) {
	case *ast.BasicLit:
		intVal, err := strconv.Atoi(et.Value)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf("pushi %d\n", byte(intVal)))
		return nil
	case *ast.BinaryExpr:
		err = s.compileExpression(et.X, buf)
		if err != nil {
			return err
		}
		err = s.compileExpression(et.Y, buf)
		if err != nil {
			return err
		}
		switch et.Op {
		case token.ADD:
			buf.WriteString("add\n")
		case token.SUB:
			buf.WriteString("sub\n")
		case token.MUL:
			buf.WriteString("mul\n")
		case token.QUO:
			buf.WriteString("div\n")
		default:
			return fmt.Errorf("unsupported token %s", et.Op.String())
		}
		return nil
	case *ast.Ident:
		loc, ok := s.identifiers[et.Name]
		if !ok {
			return fmt.Errorf("unrecognized identifier %s", et.Name)
		}
		buf.WriteString(fmt.Sprintf("push %d\n", loc))
		return nil
	case *ast.ParenExpr:
		err = s.compileExpression(et.X, buf)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unrecognized expression %T", e)
	}
}

func (s *scope) compileStatement(stmt ast.Stmt, buf *bytes.Buffer) error {
	var err error
	switch stmtType := stmt.(type) {
	case *ast.BlockStmt:
		for _, subStmt := range stmtType.List {
			err = s.compileStatement(subStmt, buf)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.ReturnStmt:
		res := stmtType.Results[0]
		err := s.compileExpression(res, buf)
		if err != nil {
			return err
		}
		buf.WriteString("pop 0\n")
		buf.WriteString("halt\n")
		return nil
	default:
		return fmt.Errorf("unsupported statement type %T", stmtType)
	}
}
