.PHONY: all windows macos
all: windows macos

windows:
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o bin/app-windows-amd64.exe main.go

macos:
	fyne package -os darwin -icon icon.png