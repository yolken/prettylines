package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
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

	// Do an initial gofmt run on the file
	var err error
	contents, err = runGoFmt(contents)
	if err != nil {
		log.Fatalf("Error formatting file: %+v", err)
	}

	for {
		log.Debugf("Starting round %d", round)

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

		fmt.Print(string(contents))

		round++
	}

	return removeAnnotations(contents)
}

func runGoFmt(contents []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	path, err := exec.LookPath("gofmt")
	if err != nil {
		log.Fatal("Cannot find gofmt in path")
	}

	cmd := exec.Cmd{
		Path:   path,
		Stdin:  bytes.NewReader(contents),
		Stdout: output,
	}

	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func formatNode(node dst.Node) {
	switch n := node.(type) {
	case dst.Decl:
		log.Debugf("Processing declaration: %+v", n)
		formatDecl(n)
	case dst.Expr:
		log.Debugf("Processing expression: %+v", n)
		formatExpr(n, false)
	case dst.Stmt:
		log.Debugf("Processing statement: %+v", n)
		formatStmt(n)
	}
}

func formatDecl(decl dst.Decl) {
	switch d := decl.(type) {
	case *dst.FuncDecl:
		if hasAnnotation(decl) {
			if d.Type != nil && d.Type.Params != nil {
				for f, field := range d.Type.Params.List {
					formatField(f, field)
				}
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
	shouldShorten := hasAnnotation(stmt)

	switch s := stmt.(type) {
	case *dst.ExprStmt:
		formatExpr(s.X, shouldShorten)
	case *dst.AssignStmt:
		for _, expr := range s.Rhs {
			formatExpr(expr, shouldShorten)
		}
	default:
		log.Debugf("Got another statement type: %+v", reflect.TypeOf(s))
	}
}

func formatExpr(expr dst.Expr, force bool) {
	shouldShorten := force || hasAnnotation(expr)

	switch e := expr.(type) {
	case *dst.CallExpr:
		for a, arg := range e.Args {
			if shouldShorten {
				if a == 0 {
					arg.Decorations().Before = dst.NewLine
				}
				arg.Decorations().After = dst.NewLine
			}
			formatExpr(arg, false)
		}
	case *dst.BinaryExpr:
		formatExpr(e.X, shouldShorten)
		formatExpr(e.Y, shouldShorten)
	default:
		log.Debugf("Got another expression type: %+v", reflect.TypeOf(e))
	}
}

func hasAnnotation(node dst.Node) bool {
	startDecorations := node.Decorations().Start.All()
	return len(startDecorations) > 0 && isAnnotation(startDecorations[len(startDecorations)-1])
}
