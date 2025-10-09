package args

var jsonHelp string = `Afbd help for json target.

Flags:
  -help  Display help.

Parameters:
  -path  Path for output files.
`

var Json json

type json struct {
	Present bool
	// Parameters
	Path string
}

func isValidFlagJson(flag string) bool {
	validFlags := map[string]bool{
		"-help": true,
	}

	if _, ok := validFlags[flag]; ok {
		return true
	}

	return false
}

func isValidParamJson(param string) bool {
	validParams := map[string]bool{
		"-path": true,
	}

	if _, ok := validParams[param]; ok {
		return true
	}

	return false
}

func setFlagJson(flag string) {}

func setParamJson(arg string) {
	switch param {
	case "-path":
		Json.Path = arg
	}
}

func postprocessJson() {
	if Json.Path == "" {
		Json.Path = Path
	}
}
