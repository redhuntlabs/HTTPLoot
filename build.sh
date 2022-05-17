export CGO_ENABLED=0
GOOS=linux GOARCH=amd64 go build -o httploot-linux64
GOOS=darwin GOARCH=amd64 go build -o httploot-darwin64
GOOS=windows GOARCH=amd64 go build -o httploot-windows64.exe
GOOS=windows GOARCH=386 go build -o httploot-windows32.exe
GOOS=freebsd GOARCH=amd64 go build -o httploot-freebsd64
GOOS=openbsd GOARCH=amd64 go build -o httploot-openbsd64
shasum -a 256 httploot-* > checksums.txt