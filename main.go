package main

import (
	"TabStop/server"
	"TabStop/ui"
)

func main() {
	go server.Start()
	ui.Run()
}
