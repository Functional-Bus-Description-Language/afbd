package args

var pythonHelp string = `Afbd help for Python target.

Flags:
  -help  Display help.

Parameters:
  -path  Path for output files.
`

var Python python

type python struct {
	Present bool
	// Parameters
	Path string
}

func isValidFlagPython(flag string) bool {
	validFlags := map[string]bool{
		"-help": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidParamPython(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func setFlagPython(flag string) {}

func setParamPython(arg string) {
	switch param {
	case "-path":
		Python.Path = arg
	}
}

func postprocessPython() {
	if Python.Path == "" {
		Python.Path = Path
	}
}
