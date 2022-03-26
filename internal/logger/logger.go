package logger

import (
	"log"

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
