package logger

import (
	"os"
	"sensible/internal/constants"
	"time"
)

func log(prefix, color, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// Preallocate and build the full message
	full := "[" + timestamp + "] " + color + prefix + constants.ColorReset + " " + msg + "\n"
	_, _ = os.Stdout.WriteString(full)
}

func Info(msg string) {
	log("INFO", constants.ColorBlue, msg)
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
