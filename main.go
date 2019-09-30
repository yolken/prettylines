package main

import (
	"io/ioutil"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	opts := &Opts{}
	args, err := flags.Parse(opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, arg := range args {
		format(arg)
	}

	log.Println("Done!")
}

func format(filePath string) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %+v", filePath, err)
	}

	result, err := decorator.Parse(contents)
	if err != nil {
		log.Fatalf("Error parsing file %s: %+v", filePath, err)
	}

	for _, decl := range result.Decls {
		formatNode(decl)
	}

	if err := decorator.Print(result); err != nil {
		log.Fatalf("Error printing result: %+v", err)
	}
}

func formatNode(node dst.Node) {
	switch n := node.(type) {
	case dst.Decl:
		log.Debugf("Processing declaration: %+v", n)
		formatDecl(n)
	case dst.Expr:
		log.Debugf("Processing expression: %+v", n)
	case dst.Stmt:
		log.Debugf("Processing statement: %+v", n)
	}
}

func formatDecl(decl dst.Decl) {
	switch d := decl.(type) {
	case *dst.FuncDecl:
		log.Debugf("Got a function declaration: %+v", d)
		if d.Type != nil && d.Type.Params != nil && len(d.Type.Params.List) > 0 {
			for f, field := range d.Type.Params.List {
				formatField(f, field)
			}
		}
	default:
		log.Debugf("Got a non-function declaration")
	}
}

func formatField(f int, field *dst.Field) {
	if f == 0 {
		field.Decorations().Before = dst.NewLine
	}
	field.Decorations().After = dst.NewLine
}
