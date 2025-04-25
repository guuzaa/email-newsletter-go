package main

import (
	"github.com/guuzaa/email-newsletter/internal"
)

func main() {
	config, err := internal.Configuration("configuration")
	if err != nil {
		panic(err)
	}
	internal.Run(&config)
}
