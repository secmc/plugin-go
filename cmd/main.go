package main

import (
	"log"

	"github.com/secmc/plugin-go/plugin"
)

func main() {
	_, err := plugin.NewPlugin("test")
	if err != nil {
		log.Fatalln(err)
	}
	select {}
}
