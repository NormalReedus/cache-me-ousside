package main

import (
	"fmt"

	"github.com/NormalReedus/lru-cache-microservice/internal/config"
)

func main() {
	fmt.Println(config.Load("./config.example.json5"))
}
