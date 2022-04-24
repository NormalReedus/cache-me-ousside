package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

//TODO: take a log file from config to output to if provided, otherwise use Stdout
//TODO: write tests for log file (check for file existence and creation in testdata/ when used, remember to delete file again)

const (
	DEFAULT_FLAGS = log.Ldate | log.Ltime | log.Lmsgprefix
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

	infoLog.SetFlags(DEFAULT_FLAGS)
	warningLog.SetFlags(DEFAULT_FLAGS)
	errorLog.SetFlags(DEFAULT_FLAGS)

	// clrInfo := color.New(color.Bold)
	// infoLog = log.New(os.Stdout, clrInfo.Sprint("ℹ️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

	// clrWarn := color.New(color.FgYellow, color.Bold)
	// warningLog = log.New(os.Stdout, clrWarn.Sprint("⚠️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

	// clrErr := color.New(color.FgRed, color.Bold)
	// errorLog = log.New(os.Stdout, clrErr.Sprint("⛔ "), log.Ldate|log.Ltime|log.Lmsgprefix)

	// Use this for CACHE [OPERATION] printing with / without color
	if logFile == nil {
		terminalMode = true
	}

	return logFile // Will be nil in terminal mode
}

// func init() {
// 	clrInfo := color.New(color.Bold)
// 	infoLog = log.New(os.Stdout, clrInfo.Sprint("ℹ️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

// 	clrWarn := color.New(color.FgYellow, color.Bold)
// 	warningLog = log.New(os.Stdout, clrWarn.Sprint("⚠️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

// 	clrErr := color.New(color.FgRed, color.Bold)
// 	errorLog = log.New(os.Stdout, clrErr.Sprint("⛔ "), log.Ldate|log.Ltime|log.Lmsgprefix)
// }

func setLogFileMode(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Panic(fmt.Errorf("Could not set log file %q, got the following error: %v", filepath, err))
	}

	infoLog.SetOutput(file)
	infoLog.SetPrefix("INFO - ")
	warningLog.SetOutput(file)
	warningLog.SetPrefix("WARN - ")
	errorLog.SetOutput(file)
	errorLog.SetPrefix("ERROR - ")

	return file
}

func setTerminalMode() {
	// infoLog := log.New(os.Stdout, clrInfo.Sprint("ℹ️ "), log.Ldate|log.Ltime|log.Lmsgprefix)
	clrInfo := color.New(color.Bold)
	infoLog.SetOutput(os.Stdout)
	infoLog.SetPrefix(clrInfo.Sprint("ℹ️ "))

	// warningLog := log.New(os.Stdout, clrWarn.Sprint("⚠️ "), log.Ldate|log.Ltime|log.Lmsgprefix)
	clrWarn := color.New(color.FgYellow, color.Bold)
	warningLog.SetOutput(os.Stdout)
	warningLog.SetPrefix(clrWarn.Sprint("⚠️ "))

	// errorLog := log.New(os.Stdout, clrErr.Sprint("⛔ "), log.Ldate|log.Ltime|log.Lmsgprefix)
	clrErr := color.New(color.FgRed, color.Bold)
	errorLog.SetOutput(os.Stdout)
	errorLog.SetPrefix(clrErr.Sprint("⛔ "))

	// return infoLog, warningLog, errorLog
}

func CacheWrite(key string) {
	msg := "CACHE WRITE => " + key

	if terminalMode {
		clr := color.New(color.FgBlue, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheRead(key string) {
	msg := "CACHE READ => " + key

	if terminalMode {
		clr := color.New(color.FgGreen, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheEvict(key string) {
	msg := "CACHE EVICT => " + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheBust(key string) {
	msg := "CACHE BUST => " + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

func CacheSkip(key string) {
	msg := "CACHE SKIP => " + key

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

// func Fatal(err error) {
// 	errorLog.Fatalln(err)
// }

func HiMom(apiUrl string, port string) {
	urlClr := color.New(color.FgHiGreen, color.Underline)
	cacheClr := color.New(color.FgBlue, color.Underline)

	fmt.Println()
	fmt.Println("Your LRU cache microservice is caching requests to your proxied API.")
	fmt.Println()
	fmt.Println("Proxied API: " + urlClr.Sprint(apiUrl))
	fmt.Println("Cache URL: " + cacheClr.Sprint("http://localhost:"+port))
}
