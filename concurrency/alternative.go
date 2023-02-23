package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func createUser(name string) {
	defer wg.Done()
	fmt.Println(name + " has been created successfully")
}

func uploadUserImage() {
	defer wg.Done()
	time.Sleep(time.Second * 3)
	fmt.Println("User image has been successfully uploaded to the cloud !!")
}

func main2() {
	go createUser("test")
	go uploadUserImage()

	wg.Add(2)
	wg.Wait()
}
