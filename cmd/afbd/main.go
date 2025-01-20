package main

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/csync"
	"github.com/Functional-Bus-Description-Language/afbd/internal/json"
	"github.com/Functional-Bus-Description-Language/afbd/internal/python"
	"github.com/Functional-Bus-Description-Language/afbd/internal/vhdlwb3"

	"fmt"
	"log"
	"os"
)

var printDebug bool = false

type Logger struct{}

func (l Logger) Write(p []byte) (int, error) {
	print := true

	if len(p) > 4 && string(p)[:5] == "debug" {
		print = printDebug
	}

	if print {
		fmt.Fprint(os.Stderr, string(p))
	}

	return len(p), nil
}

func main() {
	logger := Logger{}
	log.SetOutput(logger)
	log.SetFlags(0)

	cmdLineArgs := args.Parse()
	args.SetOutputPaths(cmdLineArgs)

	if _, ok := cmdLineArgs["global"]["--debug"]; ok {
		printDebug = true
	}

	mainName := "Main"
	if _, ok := cmdLineArgs["global"]["-main"]; ok {
		mainName = cmdLineArgs["global"]["-main"]
	}
	addTimestamp := false
	if _, ok := cmdLineArgs["global"]["-add-timestamp"]; ok {
		addTimestamp = true
	}
	bus, pkgsConsts, err := fbdl.Compile(cmdLineArgs["global"]["main"], mainName, addTimestamp)
	if err != nil {
		log.Fatalf("compile: %v", err)
	}

	if _, ok := cmdLineArgs["json"]; ok {
		json.Generate(bus, pkgsConsts, cmdLineArgs["json"])
	}

	if _, ok := cmdLineArgs["c-sync"]; ok {
		csync.Generate(bus, pkgsConsts, cmdLineArgs["c-sync"])
	}

	if _, ok := cmdLineArgs["python"]; ok {
		python.Generate(bus, pkgsConsts, cmdLineArgs["python"])
	}

	if _, ok := cmdLineArgs["vhdl-wb3"]; ok {
		vhdlwb3.Generate(bus, pkgsConsts, cmdLineArgs["vhdl-wb3"])
	}
}