package cli

import "runtime"

func Red(text string) string {
	return colorize("\033[31m", text)
}

func Green(text string) string {
	return colorize("\033[32m", text)
}

func Yellow(text string) string {
	return colorize("\033[33m", text)
}

func Blue(text string) string {
	return colorize("\033[34m", text)
}

func Purple(text string) string {
	return colorize("\033[35m", text)
}

func Cyan(text string) string {
	return colorize("\033[36m", text)
}

func Gray(text string) string {
	return colorize("\033[37m", text)
}

func White(text string) string {
	return colorize("\033[97m", text)
}

func colorize(color string, text string) string {
	if runtime.GOOS == "windows" {
		return text
	}

	return color + text + "\033[0m"
}
