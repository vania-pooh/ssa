# Simple Subscribe Application
A static page with form and an API to save emails to file.

## Building and Running
1. Install [Govendor](https://github.com/kardianos/govendor)
2. Install [go-bindata](https://github.com/jteeuwen/go-bindata)
3. Build source and run:
```
$ govendor sync
$ go-bindata -o static.go static/
$ go build
$ ./ssa -help
```