package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/fatih/color"
)

//TODO: take a log file from config to output to if provided, otherwise use Stdout

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func init() {
	clrInfo := color.New(color.FgBlue, color.Bold)
	infoLogger = log.New(os.Stdout, clrInfo.Sprint("ℹ️ INFO - "), log.Ldate|log.Ltime|log.Lmsgprefix)

	clrWarn := color.New(color.FgYellow, color.Bold)
	warningLogger = log.New(os.Stdout, clrWarn.Sprint("⚠️ WARN - "), log.Ldate|log.Ltime|log.Lmsgprefix)

	clrErr := color.New(color.FgRed, color.Bold)
	errorLogger = log.New(os.Stdout, clrErr.Sprint("⛔ ERROR - "), log.Ldate|log.Ltime|log.Lmsgprefix)
}

func CacheWrite(key string) {
	clr := color.New(color.FgBlue, color.Bold)
	infoLogger.Println(clr.Sprint("CACHE WRITE: " + key))
}

func CacheRead(key string) {
	clr := color.New(color.FgGreen, color.Bold)
	log.Println(clr.Sprint("CACHE READ: " + key))
}

func CacheEvict(key string) {
	clr := color.New(color.FgRed, color.Bold)
	log.Println(clr.Sprint("CACHE EVICT: " + key))
}

func CacheBust(key string) {
	clr := color.New(color.FgRed, color.Bold)
	log.Println(clr.Sprint("CACHE BUST: " + key))
}

func CacheSkip(key string) {
	clr := color.New(color.FgYellow, color.Bold)
	log.Println(clr.Sprint("CACHE SKIP: " + key))
}

func HiMom(conf *config.Config, port string) {
	cacheColor := color.New(color.FgBlue, color.Bold)
	urlColor := color.New(color.FgHiGreen, color.Underline)

	fmt.Print("Your ")
	cacheColor.Print("LRU cache microservice ")
	fmt.Printf("is being served on http://localhost:%v.\n", port)
	fmt.Print("All requests will be proxied to ")
	urlColor.Println(conf.ApiUrl + "\n")
}
