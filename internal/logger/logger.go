package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

//TODO: take a log file from config to output to if provided, otherwise use Stdout
//TODO: write tests for log file (check for file existence and creation in testdata/ when used, remember to delete file again)

var (
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
)

func init() {
	clrInfo := color.New(color.Bold)
	// infoLog = log.New(os.Stdout, clrInfo.Sprint("ℹ️ INFO - "), log.Ldate|log.Ltime|log.Lmsgprefix)
	infoLog = log.New(os.Stdout, clrInfo.Sprint("ℹ️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

	clrWarn := color.New(color.FgYellow, color.Bold)
	// warningLog = log.New(os.Stdout, clrWarn.Sprint("⚠️ WARN - "), log.Ldate|log.Ltime|log.Lmsgprefix)
	warningLog = log.New(os.Stdout, clrWarn.Sprint("⚠️ "), log.Ldate|log.Ltime|log.Lmsgprefix)

	clrErr := color.New(color.FgRed, color.Bold)
	// errorLog = log.New(os.Stdout, clrErr.Sprint("⛔ ERROR - "), log.Ldate|log.Ltime|log.Lmsgprefix)
	errorLog = log.New(os.Stdout, clrErr.Sprint("⛔ "), log.Ldate|log.Ltime|log.Lmsgprefix)
}

func CacheWrite(key string) {
	clr := color.New(color.FgBlue, color.Bold)
	// infoLog.Println(clr.Sprint("CACHE WRITE: " + key))
	infoLog.Println(clr.Sprint("CACHE WRITE - " + key))
}

func CacheRead(key string) {
	clr := color.New(color.FgGreen, color.Bold)
	// infoLog.Println(clr.Sprint("CACHE READ: " + key))
	infoLog.Println(clr.Sprint("CACHE READ - " + key))
}

func CacheEvict(key string) {
	clr := color.New(color.FgRed, color.Bold)
	// infoLog.Println(clr.Sprint("CACHE EVICT: " + key))
	infoLog.Println(clr.Sprint("CACHE EVICT - " + key))
}

func CacheBust(key string) {
	clr := color.New(color.FgRed, color.Bold)
	// infoLog.Println(clr.Sprint("CACHE BUST: " + key))
	infoLog.Println(clr.Sprint("CACHE BUST - " + key))
}

func CacheSkip(key string) {
	clr := color.New(color.FgYellow, color.Bold)
	// infoLog.Println(clr.Sprint("CACHE SKIP: " + key))
	infoLog.Println(clr.Sprint("CACHE SKIP - " + key))
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
