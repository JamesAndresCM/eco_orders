package logger

import (
	"log"
	"os"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

var Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func Info(v ...any) {
	Logger.Println(Blue + "INFO: " + Reset, v)
}

func Success(v ...any) {
	Logger.Println(Green + "SUCCESS: " + Reset, v)
}

func Warn(v ...any) {
	Logger.Println(Yellow + "WARN: " + Reset, v)
}

func Error(v ...any) {
	Logger.Println(Red + "ERROR: " + Reset, v)
}

