package main

import (
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"
	"github.com/Functional-Bus-Description-Language/afbd/internal/csync"
	"github.com/Functional-Bus-Description-Language/afbd/internal/json"
	"github.com/Functional-Bus-Description-Language/afbd/internal/python"
	"github.com/Functional-Bus-Description-Language/afbd/internal/vhdlapb"

	"fmt"
	"log"
	"os"
)

type Logger struct{}

func (l Logger) Write(p []byte) (int, error) {
	print := true

	if len(p) > 4 && string(p)[:5] == "debug" {
		print = args.Debug
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

	args.Parse()

	bus, pkgsConsts, err := fbdl.Compile(
		args.MainFile, args.MainBus, args.AddTimestamp,
	)
	if err != nil {
		log.Fatalf("compile: %v", err)
	}

	if args.Json.Present {
		json.Generate(bus, pkgsConsts)
	}

	if args.CSync.Present {
		csync.Generate(bus, pkgsConsts)
	}

	if args.Python.Present {
		python.Generate(bus, pkgsConsts)
	}

	if args.VhdlApb.Present {
		vhdlapb.Generate(bus, pkgsConsts)
	}
}
