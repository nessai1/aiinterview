package main

import (
	"fmt"
	"github.com/nessai1/aiinterview/internal/service"
)

func main() {
	config, err := service.FetchConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("error fetching config: %v", err))
	}

	s, err := service.NewService(config)
	if err != nil {
		panic(err)
	}

	err = s.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("error starting service: %v", err))
	}
}
