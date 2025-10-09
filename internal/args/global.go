package args

// Global arguments
var (
	// Flags
	Debug        bool
	AddTimestamp bool
	Times        bool
	// Parameters
	MainBus  string
	MainFile string
	Path     string
)

func isValidTarget(target string) bool {
	validTargets := map[string]bool{
		"c-sync":   true,
		"json":     true,
		"python":   true,
		"vhdl-apb": true,
	}

	if _, ok := validTargets[target]; ok {
		return true
	}

	return false
}

func isValidFlag(flag string, target string) bool {
	switch target {
	case "c-sync":
		return isValidFlagCSync(flag)
	case "json":
		return isValidFlagJson(flag)
	case "python":
		return isValidFlagPython(flag)
	case "vhdl-apb":
		return isValidFlagVhdlApb(flag)
	default:
		panic("should never happen")
	}
}

func isValidParam(param string, target string) bool {
	if !isValidTarget(target) {
		panic("should never happen")
	}

	switch target {
	case "c-sync":
		return isValidParamCSync(param)
	case "json":
		return isValidParamJson(param)
	case "python":
		return isValidParamPython(param)
	case "vhdl-apb":
		return isValidParamVhdlApb(param)
	}

	return false
}
