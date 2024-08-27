# Makefile

build: linux # windows macos android

linux:
	go build -o fyne-cross/bin/linux-amd64/tabStop-amd64

windows:
	fyne-cross windows -arch=amd64 --app-id io.TabStop.TabStop

macos:
	fyne-cross darwin -arch=amd64 --app-id io.TabStop.TabStop