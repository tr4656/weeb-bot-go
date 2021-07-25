package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/tr4656/weeb-bot-go/framework/dependency"
)

func main() {
	fmt.Println("Weeb bot")

	dep := dependency.New()
	var wg sync.WaitGroup
	wg.Add(3)
	defer fmt.Println("All waited done")

	go func() {
		dep.Await()
		fmt.Println("Dep ready")
		wg.Done()
	}()

	go func() {
		fmt.Println("Gonna sleep")
		time.Sleep(5000 * time.Millisecond)
		fmt.Println("Waking up, readying dep")
		dep.Ready()
		fmt.Println("Done dep ready")
		wg.Done()
	}()

	go func() {
		time.Sleep(7000 * time.Millisecond)
		fmt.Println("Starting late, waiting on dep")
		dep.Await()
		fmt.Println("Done being late")
		wg.Done()
	}()

	wg.Wait()
}
