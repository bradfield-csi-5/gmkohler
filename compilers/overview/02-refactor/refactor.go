package main

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"log"
	"os"
	"slices"
	"strings"
)

const src string = `package foo

import (
	"fmt"
	"time"
)

func baz() {
	fmt.Println("Hello, world!")
}

type A int

const b = "testing"

func bar() {
	fmt.Println(time.Now())
}`

// Moves all top-level functions to the end, sorted in alphabetical order.
// The "source file" is given as a string (rather than e.g. a filename).
func SortFunctions(src string) (string, error) {
	parsed, err := decorator.Parse(src)

	if err != nil {
		return "", err
	}
	slices.SortStableFunc(parsed.Decls, func(a dst.Decl, b dst.Decl) int {
		switch aType := a.(type) {
		case *dst.FuncDecl:
			// if both are funcs then sort by name
			if bFunc, ok := b.(*dst.FuncDecl); ok {
				return strings.Compare(aType.Name.Name, bFunc.Name.Name)
			}
			return 1 // otherwise A is a func and B is not so funcs go at bottom
		default:
			switch b.(type) {
			case *dst.FuncDecl: // funcs go at bottom
				return -1
			default:
				return 0
			}
		}
	})
	var buf strings.Builder
	if err = decorator.Fprint(&buf, parsed); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	f, err := decorator.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	// Print AST
	err = dst.Fprint(os.Stdout, f, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Convert AST back to source
	err = decorator.Print(f)
	if err != nil {
		log.Fatal(err)
	}
}
