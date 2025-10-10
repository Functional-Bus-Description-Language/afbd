package args

var jsonHelp string = `Afbd help for json target.

Flags:
  -help  Display help.

Parameters:
  -path        Path for output files.
  -const-name  Name of the packages constants output json file.
               The default name is 'const.json'.
  -reg-name    Name of the registerification output json file.
               The default name is 'reg.json'.
`

var Json json

type json struct {
	Present bool
	// Parameters
	Path      string
	ConstName string
	RegName   string
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
		"-path": true, "-const-name": true, "-reg-name": true,
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
	case "-const-name":
		Json.ConstName = arg
	case "-reg-name":
		Json.RegName = arg
	}
}

func postprocessJson() {
	if Json.Path == "" {
		Json.Path = Path
	}
	if Json.ConstName == "" {
		Json.ConstName = "const.json"
	}
	if Json.RegName == "" {
		Json.RegName = "reg.json"
	}
}
