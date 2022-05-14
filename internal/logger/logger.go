package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

const (
	defaultFlags    = log.Ldate | log.Ltime | log.Lmsgprefix // defaultFlags to set logs to use date, time, and sets prefix to after the date and time.
	prefixSeparator = " => "                                 // Separator between cache operation indicator and the requested route.
	// Prefix for the type of log message.
	infoPrefix  = "INFO - "
	warnPrefix  = "WARN - "
	errorPrefix = "ERROR - "
)

var (
	infoLog    = new(log.Logger)
	warningLog = new(log.Logger)
	errorLog   = new(log.Logger)

	terminalMode bool
)

// Initialize configures the logger service to use a log file from logFilepath or run in terminal mode.
// Returns a reference to the open log file, if a log file is specified.
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

// setLogFileMode configures the logger to use a file at filepath.
// Returns a reference to the open log file.
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

// setTerminalMode configures the logger to run in terminal mode.
// This means using emojis and colors.
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

// CacheWrite will log a formatted message for a cache write operation to key
// with correct colors and cache operation indicator.
func CacheWrite(key string) {
	msg := "CACHE WRITE" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgBlue, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

// CacheRead will log a formatted message for a cache read operation to key
// with correct colors and cache operation indicator.
func CacheRead(key string) {
	msg := "CACHE READ" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgGreen, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

// CacheEvict will log a formatted message for a cache evict operation to key
// with correct colors and cache operation indicator.
func CacheEvict(key string) {
	msg := "CACHE EVICT" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

// CacheBust will log a formatted message for a cache bust operation to key
// with correct colors and cache operation indicator.
func CacheBust(key string) {
	msg := "CACHE BUST" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgRed, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

// CacheSkip will log a formatted message for a cache skip operation to key
// with correct colors and cache operation indicator.
func CacheSkip(key string) {
	msg := "CACHE SKIP" + prefixSeparator + key

	if terminalMode {
		clr := color.New(color.FgYellow, color.Bold)
		msg = clr.Sprint(msg)
	}

	infoLog.Println(msg)
}

// Info will log msg with the infoPrefix and correct icon.
func Info(msg string) {
	infoLog.Println(msg)
}

// Warn will log msg with the warnPrefix and correct icon.
func Warn(msg string) {
	warningLog.Println(msg)
}

// Error will log err with the errorPrefix and correct icon.
func Error(err error) {
	errorLog.Println(err)
}

// Panic will log err with the errorPrefix and correct icon
// as well as stop execution.
// This is only used for errors during setup.
func Panic(err error) {
	errorLog.Panicln(err)
}

// HiMom will display a startup message with a presentation of used configuration.
func HiMom(confString string, url string) {
	urlClr := color.New(color.FgBlue, color.Underline)

	myFigure := figure.NewColorFigure("Hello World", "", "green", true)
	myFigure.Print()

	fmt.Printf("📦 You LRU cache microservice is running on %s with the following configuration:\n", urlClr.Sprint(url))

	fmt.Println(confString)
}
