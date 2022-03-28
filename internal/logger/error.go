package logger

import (
	"log"

	"github.com/fatih/color"
)

func Error(err error) {
	clr := color.New(color.BgRed, color.FgWhite, color.Bold)
	log.Println(clr.Sprint("CACHE WRITE: " + err.Error()))
}
