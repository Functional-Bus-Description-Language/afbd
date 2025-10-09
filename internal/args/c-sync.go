package args

var csyncHelp string = `Afbd help for C-Sync target.
C-Sync target is a C language target with synchronous (blocking) interface
functions.

Flags:
  -help        Display help.
  -no-asserts  Do not include asserts. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var CSync cSync

type cSync struct {
	Present bool
	// Flags
	NoAsserts bool
	// Parameters
	Path string
}

func isValidFlagCSync(flag string) bool {
	validFlags := map[string]bool{
		"-help":       true,
		"-no-asserts": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidParamCSync(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func setFlagCSync(flag string) {
	switch flag {
	case "-no-asserts":
		CSync.NoAsserts = true
	}
}

func setParamCSync(arg string) {
	switch param {
	case "-path":
		CSync.Path = arg
	}
}

func postprocessCSync() {
	if CSync.Path == "" {
		CSync.Path = Path
	}
}
