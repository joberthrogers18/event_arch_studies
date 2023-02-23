package main

import (
	"fmt"
)

func createUser(name string, done chan<- bool) {
	fmt.Println(name + " has been created successfully")
	done <- true
}

func uploadUserImage(done chan<- bool) {
	fmt.Println("User image has been successfully uploaded to the cloud !!")
	done <- true
}

func main() {
	done := make(chan bool)

	go uploadUserImage(done)
	fmt.Println("Image uploaded: ", <-done)

	go createUser("test", done)
	fmt.Println("User created: ", <-done)
}
