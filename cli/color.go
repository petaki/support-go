package cli

import "runtime"

// Red function.
func Red(text string) string {
	return colorize("\033[31m", text)
}

// Green function.
func Green(text string) string {
	return colorize("\033[32m", text)
}

// Yellow function.
func Yellow(text string) string {
	return colorize("\033[33m", text)
}

// Blue function.
func Blue(text string) string {
	return colorize("\033[34m", text)
}

// Purple function.
func Purple(text string) string {
	return colorize("\033[35m", text)
}

// Cyan function.
func Cyan(text string) string {
	return colorize("\033[36m", text)
}

// Gray function.
func Gray(text string) string {
	return colorize("\033[37m", text)
}

// White function.
func White(text string) string {
	return colorize("\033[97m", text)
}

func colorize(color string, text string) string {
	if runtime.GOOS == "windows" {
		return text
	}

	return color + text + "\033[0m"
}
