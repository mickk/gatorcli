package main

import (
	"fmt"

	"github.com/mickk/gatorcli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
	}

	cfg.SetUser("mike")
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
	}

	fmt.Printf("Current Config: %+v\n", cfg)
}
