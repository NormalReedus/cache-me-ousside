package logger

import (
	"fmt"
	"log"

	"github.com/NormalReedus/cache-me-ousside/internal/config"
	"github.com/fatih/color"
)

func CacheWrite(key string) {
	clr := color.New(color.FgBlue, color.Bold)
	log.Printf("%v\n", clr.Sprint("CACHE WRITE: "+key))
}

func CacheRead(key string) {
	clr := color.New(color.FgGreen, color.Bold)
	log.Printf("%v\n", clr.Sprint("CACHE READ: "+key))
}

func CacheEvict(key string) {
	clr := color.New(color.FgRed, color.Bold)
	log.Printf("%v\n", clr.Sprint("CACHE EVICT: "+key))
}

func CacheSkip(key string) {
	clr := color.New(color.FgYellow, color.Bold)
	log.Printf("%v\n", clr.Sprint("CACHE SKIP: "+key))
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
