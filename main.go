package main

import (
	"sync"
	"tech-test-2-MNC/api"
	"tech-test-2-MNC/config"
)

func main() {
	//Init Config
	config.SetConfig()

	var wg sync.WaitGroup
	wg.Add(1)

	//Start HTTP / REST Server
	go api.StartHttpServer()

	//Start Consumer
	nsqClient := api.NewConsumer()
	go nsqClient.RegisterAll().Run()

	wg.Wait()
}
