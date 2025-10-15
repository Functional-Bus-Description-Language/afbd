package python

import (
	_ "embed"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

var busWidth int64

//go:embed templates/afbd.py
var pythonTmplStr string
var pythonTmpl = template.Must(template.New("Python module").Parse(pythonTmplStr))

type pythonFormatters struct {}

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package) {
	busWidth = bus.Width

	err := os.MkdirAll(args.Python.Path, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}

	f, err := os.Create(path.Join(args.Python.Path, "afbd.py"))
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}

	fmts := pythonFormatters{}

	err = pythonTmpl.Execute(f, fmts)
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("generate Python: %v", err)
	}
}

var indent string

func increaseIndent(val int) {
	// NOTE: Inefficient implementation.
	for range val {
		indent += "    "
	}
}

func decreaseIndent(val int) {
	indent = indent[:len(indent)-val*4]
}
