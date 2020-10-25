package logger

import (
	"cs425_mp2/config"
	"fmt"
	"log"
	"os"
)

var (
	// InfoLogger : logging for info messages
	InfoLogger *log.Logger

	// WarningLogger : logging for warning messages
	WarningLogger *log.Logger

	// ErrorLogger : logging for error messages
	ErrorLogger *log.Logger

	// ErrorLogger : logging for error messages
	DebugLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "[WARNING]", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile)
}

func PrintToConsole(args ...interface{}) {
	fmt.Print(blue(fmt.Sprintln(args...)))
}

func PrintInfo(args ...interface{}) {
	InfoLogger.Println(args)

	fmt.Print(green("[INFO] "))
	fmt.Print(green(fmt.Sprintln(args...)))
}

func PrintWarning(args ...interface{}) {
	WarningLogger.Println(args)

	fmt.Print(yellow("[WARN] "))
	fmt.Print(yellow(fmt.Sprintln(args...)))
}

func PrintError(args ...interface{}) {
	ErrorLogger.Println(args)

	fmt.Print(red("[ERROR]"))
	fmt.Print(red(fmt.Sprintln(args...)))
}

func PrintDebug(args ...interface{}) {
	ErrorLogger.Println(args)
	if config.DebugMode {
		fmt.Print(yellow("[DEBUG] "))
		fmt.Println(args...)
	}
}