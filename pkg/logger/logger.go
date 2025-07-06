package logger

import (
	"os"
	"sensible/internal/constants"
	"strings"
)

func Info(msg ...string) {
	log("INFO", constants.ColorBlue, msg...)
}

func Success(msg ...string) {
	log("SUCCESS", constants.ColorGreen, msg...)
}

func Custom(prefix, color string, msg ...string) {
	log(prefix, color, msg...)
}

func Warn(msg ...string) {
	log("WARN", constants.ColorYellow, msg...)
}

func Error(msg ...string) {
	log("ERROR", constants.ColorRed, msg...)
}

func Plain(msg string) {
	_, _ = os.Stdout.WriteString(msg) // ignore error, as this is a simple print function
}

// helper func ...
func log(prefix, color string, msg ...string) {
	m := strings.Join(msg, " ")
	_, _ = os.Stdout.WriteString(color + "[" + prefix + "]" + constants.ColorReset + " " + m + "\n")
}
