package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

const (
	defaultFlags    = log.Ldate | log.Ltime | log.Lmsgprefix
	prefixSeparator = " => "
	infoPrefix      = "INFO - "
	warnPrefix      = "WARN - "
	errorPrefix     = "ERROR - "
)

var (
	infoLog    = new(log.Logger)
	warningLog = new(log.Logger)
	errorLog   = new(log.Logger)

	terminalMode bool
)

func Initialize(logFilepath string) *os.File {
	var logFile *os.File // will only be populated if a logfile path is provided

	// Output to a logfile (will use text instead of emojis and no colors)
	if logFilepath != "" {
		logFile = setLogFileMode(logFilepath)

	} else {
		// Only show output in terminal (stdout) (will use emojis and colors).
		setTerminalMode()
	}

	infoLog.SetFlags(defaultFlags)
	warningLog.SetFlags(defaultFlags)
	errorLog.SetFlags(defaultFlags)

	// Use this for CACHE [OPERATION] printing with / without color
	if logFile == nil {
		terminalMode = true
	}

	return logFile // Will be nil in terminal mode
}

func setLogFileMode(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Panic(fmt.Errorf("Could not set log file %q, got the following error: %v", filepath, err))
	}

	infoLog.SetOutput(file)
	infoLog.SetPrefix(infoPrefix)
	warningLog.SetOutput(file)
	warningLog.SetPrefix(warnPrefix)
	errorLog.SetOutput(file)
	errorLog.SetPrefix(errorPrefix)

	return file
}

func setTerminalMode() {
	clrInfo := color.New(color.Bold)
	infoLog.SetOutput(os.Stdout)
	infoLog.SetPrefix(clrInfo.Sprint("ℹ️ "))

	clrWarn := color.New(color.FgYellow, color.Bold)
	warningLog.SetOutput(os.Stdout)
	warningLog.SetPrefix(clrWarn.Sprint("⚠️ "))

	clrErr := color.New(color.FgRed, color.Bold)
	errorLog.SetOutput(os.Stdout)
	errorLog.SetPrefix(clrErr.Sprint("⛔ "))
}

func CacheWrite(key string) {
	msg := "CACHE WRITE" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgBlue, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheRead(key string) {
	msg := "CACHE READ" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgGreen, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheEvict(key string) {
	msg := "CACHE EVICT" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheBust(key string) {
	msg := "CACHE BUST" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheSkip(key string) {
	msg := "CACHE SKIP" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgYellow, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func Info(msg string) {
	infoLog.Println(msg)
}

func Warn(msg string) {
	warningLog.Println(msg)
}

func Error(err error) {
	errorLog.Println(err)
}

func Panic(err error) {
	errorLog.Panicln(err)
}

func HiMom(confString string, url string) {
	urlClr := color.New(color.FgBlue, color.Underline)

	fmt.Printf("You LRU cache microservice is running on %s with the following configuration:\n", urlClr.Sprint(url))

	fmt.Println(confString)
}
