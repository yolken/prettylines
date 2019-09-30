package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	log "github.com/sirupsen/logrus"
)

func shorten(filePath string) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %+v", filePath, err)
	}

	contents = processContents(contents)
	fmt.Print(string(contents))
}

func processContents(contents []byte) []byte {
	round := 0

	for {
		log.Debugf("Starting round %d", round)

		var err error
		contents, err = formatFile(contents)
		if err != nil {
			log.Fatalf("Error formatting file: %+v", err)
		}

		var linesToShorten int
		contents, linesToShorten = annotateLongLines(contents)
		if linesToShorten == 0 {
			log.Debugf("No more lines to shorten, breaking")
			break
		}

		result, err := decorator.Parse(contents)
		if err != nil {
			log.Fatalf("Error parsing file: %+v", err)
		}

		for _, decl := range result.Decls {
			formatNode(decl)
		}

		output := bytes.NewBuffer([]byte{})
		log.Printf("Generating new outputs")
		err = decorator.Fprint(output, result)
		if err != nil {
			log.Fatal(err)
		}
		contents = output.Bytes()
		round++
	}

	return removeAnnotations(contents)
}

func formatNode(node dst.Node) {
	switch n := node.(type) {
	case dst.Decl:
		log.Debugf("Processing declaration: %+v", n)
		formatDecl(n)
	case dst.Expr:
		log.Debugf("Processing expression: %+v", n)
		formatExpr(n)
	case dst.Stmt:
		log.Debugf("Processing statement: %+v", n)
		formatStmt(n)
	}
}

func formatDecl(decl dst.Decl) {
	switch d := decl.(type) {
	case *dst.FuncDecl:
		log.Debugf("Got a function declaration: %+v", d)
		if d.Type != nil && d.Type.Params != nil {
			for f, field := range d.Type.Params.List {
				formatField(f, field)
			}
		}

		if d.Body != nil && d.Body.List != nil && len(d.Body.List) > 0 {
			for _, stmt := range d.Body.List {
				formatStmt(stmt)
			}
		}
	default:
		log.Debugf("Got another type of declaration: %+v", reflect.TypeOf(d))
	}
}

func formatField(f int, field *dst.Field) {
	if f == 0 {
		field.Decorations().Before = dst.NewLine
	}
	field.Decorations().After = dst.NewLine
}

func formatStmt(stmt dst.Stmt) {
	switch s := stmt.(type) {
	case *dst.ExprStmt:
		formatExpr(s.X)
	default:
		log.Debugf("Got another statement type: %+v", reflect.TypeOf(s))
	}
}

func formatExpr(expr dst.Expr) {
	switch e := expr.(type) {
	case *dst.CallExpr:
		for a, arg := range e.Args {
			if a == 0 {
				arg.Decorations().Before = dst.NewLine
			}
			arg.Decorations().After = dst.NewLine
		}
	default:
		log.Debugf("Got another expression type: %+v", reflect.TypeOf(e))
	}
}
