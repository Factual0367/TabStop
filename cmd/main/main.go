package main

import (
	"TabStop/server"
	"TabStop/ui"
)

func main() {
	go server.Start() // Start server in background
	ui.Run()          // Start UI
}
