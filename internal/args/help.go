package args

import (
	"fmt"
	"os"
)

var helpMsg string = `Functional Bus Description Language compiler back-end for Advanced Microcontroller Bus Architecture 5 (AMBA5) specifications.
Version: %s

Supported targets:
  - c-sync    C target with synchronous (blocking) interface functions,
  - json      JSON target,
  - python    Python target,
  - vhdl-apb  VHDL target for APB.
To check valid flags and parameters for a given target type: 'afbd {target} -help'.

Usage:
  afbd [global flag or parameter] [{{target}} [target flag or parameter] ...] ... path/to/fbd/file/with/main/bus

  At least one target must be specified. The last argument is always a path
  to the fbd file containing a definition of the main bus, unless it is
  '-help' or '-version.'

Flags:
  -help           Display help.
  -version        Display version.
  -debug          Print debug messages.
  -add-timestamp  Add bus generation bus timestamp.
  -times          Print compile and generate times. Not yet implemented.

Parameters:
  -main          Name of the main bus. Useful for testbenches.
  -path          Path for target directories with output files.
                 The default is 'afbd' directory in the current working directory.
`

func printHelp() {
	fmt.Printf(helpMsg, Version)
	os.Exit(0)
}

func printTargetHelp(target string) {
	switch target {
	case "c-sync":
		fmt.Print(csyncHelpMsg)
	case "json":
		fmt.Print(jsonHelpMsg)
	case "python":
		fmt.Print(pythonHelpMsg)
	case "vhdl-apb":
		fmt.Print(vhdlAPBHelpMsg)
	default:
		panic("should never happen")
	}

	os.Exit(0)
}

var csyncHelpMsg string = `Afbd help for C-Sync target.
C-Sync target is a C language target with synchronous (blocking) interface
functions.

Flags:
  -help        Display help.
  -no-asserts  Do not include asserts. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var pythonHelpMsg string = `Afbd help for Python target.

Flags:
  -help  Display help.

Parameters:
  -path  Path for output files.
`

var vhdlAPBHelpMsg string = `Afbd help for vhdl-apb target.

Flags:
  -help   Display help.
  -no-psl Do not include PSL assertions. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var jsonHelpMsg string = `Afbd help for json target.

Flags:
  -help  Display help.

Parameters:
  -path  Path for output files.
`
