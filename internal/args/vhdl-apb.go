package args

var vhdlApbHelp string = `Afbd help for vhdl-apb target.

Flags:
  -help        Display help.
  -shared-bus  Use Shared Bus as internal interconnect instead of NxM Crossbar.
               If the number of requesters or completers is greater than 1,
               then a classic NxM Crossbar is used as an internal interconnect
               by default. If this flag is present, a Shared Bus will be used
               instead.

Parameters:
  -path  Path for output files.
`

var VhdlApb vhdlApb

type vhdlApb struct {
	Present bool
	// Flags
	SharedBus bool
	// Parameters
	Path string
}

func isValidFlagVhdlApb(flag string) bool {
	validFlags := map[string]bool{
		"-help":       true,
		"-shared-bus": true,
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
	case "-shared-bus":
		VhdlApb.SharedBus = true
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
