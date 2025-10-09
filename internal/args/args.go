// Custom package for command line arguments parsing.
package args

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const Version string = "0.0.0"

// Parser state variables
var (
	target    string // Target currently being parsed
	param     string // Parameter currently being parsed
	expectArg bool   // True if next argument must be parameter value
)

func printVersion() {
	fmt.Println(Version)
	os.Exit(0)
}

func Parse() {
	if len(os.Args) == 1 {
		printHelp()
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]
		switch arg {
		case "-help":
			printHelp()
		case "-version":
			printVersion()
		default:
			if isValidTarget(arg) {
				log.Fatalf("missing main file, check 'afbd -help'")
			} else {
				log.Fatalf("'%s' is not valid target, check 'afbd -help'", arg)
			}
		}
	}

	for i, arg := range os.Args[1:] {
		if i == len(os.Args)-2 {
			if !strings.HasPrefix(arg, "-") {
				break
			}
		}

		// Parse global flags and arguments
		if target == "" {
			if arg == "-help" {
				printHelp()
			} else if arg == "-version" {
				printVersion()
			} else if expectArg {
				switch param {
				case "-main":
					MainBus = arg
				case "-path":
					Path = arg
				}
				expectArg = false
			} else if arg == "-add-timestamp" {
				AddTimestamp = true
			} else if arg == "-times" {
				Times = true
			} else if arg == "-debug" {
				Debug = true
			} else if arg == "-main" || arg == "-path" {
				param = arg
				expectArg = true
			} else if !strings.HasPrefix(arg, "-") {
				if !isValidTarget(arg) {
					log.Fatalf("'%s' is not valid target", arg)
				}
				setTarget(arg)
			}

			continue
		}

		if expectArg {
			setParam(arg)
			expectArg = false
		} else if isValidTarget(arg) {
			setTarget(arg)
		} else if !isValidParam(arg, target) && !isValidFlag(arg, target) && !expectArg {
			log.Fatalf(
				"'%s' is not valid flag or parameter for '%s' target, "+
					"run 'afbd %[2]s -help' to see valid flags and parameters",
				arg, target,
			)
		} else if arg == "-help" {
			printTargetHelp(target)
		} else if isValidFlag(arg, target) {
			setFlag(arg)
		} else if isValidParam(arg, target) {
			param = arg
			expectArg = true
		}
	}

	if expectArg {
		log.Fatalf(
			"missing argument for '%s' parameter, target '%s'",
			param, target,
		)
	}

	MainFile = os.Args[len(os.Args)-1]

	// Default values handling.
	if Path == "" {
		Path = "afbd"
	}
	if MainBus == "" {
		MainBus = "Main"
	}

	if !CSync.Present && !Json.Present && !Python.Present && !VhdlApb.Present {
		fmt.Println("no target specified, run 'afbd -help' to check valid targets")
		os.Exit(1)
	}

	postprocessCSync()
	postprocessJson()
	postprocessPython()
	postprocessVhdlApb()
}

func setTarget(t string) {
	switch t {
	case "c-sync":
		CSync.Present = true
	case "json":
		Json.Present = true
	case "python":
		Python.Present = true
	case "vhdl-apb":
		VhdlApb.Present = true
	}

	target = t
}

func setFlag(flag string) {
	switch target {
	case "c-sync":
		setFlagCSync(flag)
	case "json":
		setFlagJson(flag)
	case "python":
		setFlagPython(flag)
	case "vhdl-apb":
		setFlagVhdlApb(flag)
	}
}

func setParam(arg string) {
	switch target {
	case "c-sync":
		setParamCSync(arg)
	case "json":
		setParamJson(arg)
	case "python":
		setParamPython(arg)
	case "vhdl-apb":
		setParamVhdlApb(arg)
	}
}
