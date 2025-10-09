package args

var vhdlApbHelp string = `Afbd help for vhdl-apb target.

Flags:
  -help   Display help.
  -no-psl Do not include PSL assertions. Not yet implemented.

Parameters:
  -path  Path for output files.
`

var VhdlApb vhdlApb

type vhdlApb struct {
	Present bool
	// Flags
	NoPsl bool
	// Parameters
	Path string
}

func isValidFlagVhdlApb(flag string) bool {
	validFlags := map[string]bool{
		"-help":   true,
		"-no-psl": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidParamVhdlApb(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func setFlagVhdlApb(flag string) {
	switch flag {
	case "-no-psl":
		VhdlApb.NoPsl = true
	}
}

func setParamVhdlApb(arg string) {
	switch param {
	case "-path":
		VhdlApb.Path = arg
	}
}

func postprocessVhdlApb() {
	if VhdlApb.Path == "" {
		VhdlApb.Path = Path
	}
}
