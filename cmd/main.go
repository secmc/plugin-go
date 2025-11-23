package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/secmc/plugin-go/plugin"
)

func main() {
	fmt.Println(os.Getenv("HOST"))

	p, err := plugin.NewPlugin("test")
	if err != nil {
		log.Fatalln()
	}

	_ = p
}
