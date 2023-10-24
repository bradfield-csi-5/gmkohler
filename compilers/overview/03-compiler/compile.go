package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

const (
	templateJeqz  = "jeqz %s\n"
	templateJump  = "jump %s\n"
	templateLabel = "label %s\n"
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
	var nextLoc addr = 1
	params := make(map[string]addr)
	for _, p := range node.Type.Params.List {
		for _, n := range p.Names {
			params[n.Name] = nextLoc
			nextLoc++
		}
	}
	s := scope{
		addrsByVarName:    params,
		nextAvailableAddr: nextLoc,
	}
	err := s.compileStatement(node.Body, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type addr byte
type scope struct {
	addrsByVarName    map[string]addr
	nextAvailableAddr addr
}

func (s *scope) declareVar(name string) error {
	if _, exists := s.addrsByVarName[name]; exists {
		return fmt.Errorf("variable %s already exists", name)
	}
	s.addrsByVarName[name] = s.nextAvailableAddr
	s.nextAvailableAddr++
	return nil
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
		case token.EQL:
			buf.WriteString("eq\n")
		case token.GTR:
			buf.WriteString("gt\n")
		case token.LSS:
			buf.WriteString("lt\n")
		case token.NEQ:
			buf.WriteString("neq\n")
		case token.LEQ:
			buf.WriteString("leq\n")
		case token.GEQ:
			buf.WriteString("geq\n")
		default:
			return fmt.Errorf("unsupported token %s", et.Op.String())
		}
		return nil
	case *ast.Ident:
		loc, ok := s.addrsByVarName[et.Name]
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
	switch sType := stmt.(type) {
	case *ast.AssignStmt:
		// expect basic assignments (e.g. not dual assignments)
		lhsId, ok := sType.Lhs[0].(*ast.Ident)
		if !ok {
			return fmt.Errorf(
				"only assignments to identifiers are allowed, got %T instead",
				sType.Lhs[0],
			)
		}
		lhsAddr, ok := s.addrsByVarName[lhsId.Name]
		if !ok {
			return fmt.Errorf("unrecognized identifier %s", lhsId.Name)
		}
		err = s.compileExpression(sType.Rhs[0], buf)
		if err != nil {
			return err
		}
		/**
		 * pop what was pushed (by evaluating Rhs) from stack to the mem addr of
		 * the Lhs identifier
		 */
		buf.WriteString(fmt.Sprintf("pop %d\n", lhsAddr))
		return nil
	case *ast.BlockStmt:
		for _, subStmt := range sType.List {
			err = s.compileStatement(subStmt, buf)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.DeclStmt:
		genDecl, ok := sType.Decl.(*ast.GenDecl)
		if !ok {
			return fmt.Errorf(
				"only general declarations are supported.  received %T",
				sType,
			)
		}
		valueSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
		if !ok {
			return fmt.Errorf(
				"only val declarations are supported.  received %T",
				genDecl.Specs[0],
			)
		}
		err = s.declareVar(valueSpec.Names[0].Name)
		if err != nil {
			return err
		}
		return nil
	case *ast.ForStmt:
		// does not handle init/increment
		labelBefore := fmt.Sprintf("for-%d-before", sType.For)
		labelAfter := fmt.Sprintf("for-%d-after", sType.For)
		// we'll jump back here based on condition
		buf.WriteString(fmt.Sprintf(templateLabel, labelBefore))
		err = s.compileExpression(sType.Cond, buf)
		if err != nil {
			return err
		}
		// skip body if condition is false
		buf.WriteString(fmt.Sprintf(templateJeqz, labelAfter))
		// evaluate body
		err = s.compileStatement(sType.Body, buf)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf(templateJump, labelBefore))
		buf.WriteString(fmt.Sprintf(templateLabel, labelAfter))

		return nil
	case *ast.IfStmt:
		// requires labels for branching
		labelElse := fmt.Sprintf("if-%d-else", sType.If)
		labelAfter := fmt.Sprintf("if-%d-after", sType.If)
		// evaluate condition
		err = s.compileExpression(sType.Cond, buf)
		if err != nil {
			return err
		}
		jumpLabel := labelAfter
		if sType.Else != nil {
			jumpLabel = labelElse
		}
		// we either jump to the else or after in the false case; otherwise
		// we proceed with the body:
		buf.WriteString(fmt.Sprintf(templateJeqz, jumpLabel))
		err = s.compileStatement(sType.Body, buf)
		if err != nil {
			return err
		}
		if sType.Else != nil {
			// if we have to write the else statement, jump over it in the
			// body statement that we just wrote
			buf.WriteString(fmt.Sprintf(templateJump, labelAfter))
			buf.WriteString(fmt.Sprintf(templateLabel, labelElse))
			err = s.compileStatement(sType.Else, buf)
			if err != nil {
				return err
			}
		}
		buf.WriteString(fmt.Sprintf(templateLabel, labelAfter))
		return nil
	case *ast.ReturnStmt:
		res := sType.Results[0]
		err := s.compileExpression(res, buf)
		if err != nil {
			return err
		}
		buf.WriteString("pop 0\n")
		buf.WriteString("halt\n")
		return nil
	default:
		return fmt.Errorf("unsupported statement type %T", sType)
	}
}
