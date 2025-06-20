package logger

import (
	"os"
	"sensible/internal/constants"
)

func log(prefix, color, msg string) {
	_, _ = os.Stdout.WriteString(color + "[" + prefix + "]" + constants.ColorReset + " " + msg + "\n")
}

func Info(msg string) {
	log("INFO", constants.ColorBlue, msg)
}

func Success(msg string) {
	log("SUCCESS", constants.ColorGreen, msg)
}

func Custom(prefix, color, msg string) {
	log(prefix, color, msg)
}

func Warn(msg string) {
	log("WARN", constants.ColorYellow, msg)
}

func Error(msg string) {
	log("ERROR", constants.ColorRed, msg)
}

func Plain(msg string) {
	_, _ = os.Stdout.WriteString(msg + "\n")
}
