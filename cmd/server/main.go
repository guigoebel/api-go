package main

import "github.com/guigoebel/api-go/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
