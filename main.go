package main

import "fmt"

func main() {
	config := getConfig()
	fmt.Println(config.Auth.User)
}
