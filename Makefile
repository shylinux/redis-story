export GOPROXY=https://goproxy.cn,direct
export GOPRIVATE=shylinux.com,github.com
export CGO_ENABLED=0

all:
	@echo && date
	go build -v -o bin/ice.bin src/main.go && chmod u+x bin/ice.bin && ./bin/ice.sh restart
