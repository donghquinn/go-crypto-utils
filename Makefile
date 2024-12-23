all: windows macos

windows:
    CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 GOARCH=amd64 go build -o myapp.exe

macos:
    GOOS=darwin GOARCH=amd64 go build -o myapp