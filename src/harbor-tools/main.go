package main

import "harbor-tools/harbor-tools/di"

func main() {
	configPath := "/finup/go-projects/harbor-tools/configs/"
	server, _ :=di.InitServer(configPath)
	server()
}

