# Makefile

build: linux # windows macos android

linux:
	fyne-cross linux -arch=amd64 --app-id io.TabStop.TabStop

windows:
	fyne-cross windows -arch=amd64 --app-id io.TabStop.TabStop

macos:
	fyne-cross darwin -arch=amd64 --app-id io.TabStop.TabStop

android:
	fyne-cross android -arch=arm64,arm --app-id io.TabStop.TabStop
