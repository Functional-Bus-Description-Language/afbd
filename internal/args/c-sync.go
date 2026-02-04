package args

var csyncHelp string = `Afbd help for C-Sync target.
C-Sync target is a C language target with synchronous (blocking) interface
functions.

Flags:
  -help         Display help.
  -linux-mmap-iface
                Generate Linux memory-mapped IO interface implementation.
  -no-asserts   Do not include asserts. Not yet implemented.
  -offset-addr  Use block inner offset addresses instead of global addresses.

Parameters:
  -path  Path for output files.
`

var CSync cSync

type cSync struct {
	Present bool
	// Flags
	LinuxMmapIface bool
	NoAsserts      bool
	OffsetAddr     bool
	// Parameters
	Path string
}

func isValidFlagCSync(flag string) bool {
	validFlags := map[string]bool{
		"-help":             true,
		"-linux-mmap-iface": true,
		"-no-asserts":       true,
		"-offset-addr":      true,
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
	case "-linux-mmap-iface":
		CSync.LinuxMmapIface = true
	case "-no-asserts":
		CSync.NoAsserts = true
	case "-offset-addr":
		CSync.OffsetAddr = true
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
