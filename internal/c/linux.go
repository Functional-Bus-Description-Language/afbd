package c

// WidthToLinuxReadFunc returns name of the Linux read function to be used inside driver.
// In case of unsupported width, "XXX" string is returned instead of panicking.
func WidthToLinuxReadFunc(width int64) string {
	switch width {
	case 8:
		return "readb"
	case 16:
		return "readw"
	case 32:
		return "readl"
	case 64:
		return "readq"
	}

	return "XXX"
}

// WidthToLinuxWriteFunc returns name of the Linux write function to be used inside driver.
// In case of unsupported width, "XXX" string is returned instead of panicking.
func WidthToLinuxWriteFunc(width int64) string {
	switch width {
	case 8:
		return "writeb"
	case 16:
		return "writew"
	case 32:
		return "writel"
	case 64:
		return "writeq"
	}

	return "XXX"
}
