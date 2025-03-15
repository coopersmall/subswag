package main

import (
	"sync"

	"github.com/coopersmall/subswag/env"
	"github.com/coopersmall/subswag/http/routers/api"
	"github.com/coopersmall/subswag/streams/subscribers"
	"github.com/joho/godotenv"
)

func main() {
	run()
}

func run() {
	stop := make(chan struct{})
	closer := start()
	<-stop
	err := closer()
	if err != nil {
		panic(err)
	}
}

func start() func() error {
	godotenv.Load()

	vars := env.MustGetEnvVars(env.Opts...)
	env := env.MustGetEnv(vars)

	var stopRouterChan chan func() error
	defer close(stopRouterChan)

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()
		stopRouterChan <- api.NewAPIRouter(env).MustStart(vars)
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		s, close := env.GetSubscribers()
		defer close()
		subscribers.StartSubscribers(env, s)
	}()

	wg.Wait()

	stopRouter := <-stopRouterChan
	stopFunc := func() error {
		err := stopRouter()
		if err != nil {
			return err
		}
		return nil
	}
	return stopFunc
}
