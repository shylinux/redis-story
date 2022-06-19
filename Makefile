binarys = bin/ice.bin

all: def
	@echo && date
	go build -v -o ${binarys} src/main.go src/version.go src/binpack.go && ./${binarys} forever restart &>/dev/null

def:
	@ [ -f src/version.go ] || echo "package main" > src/version.go
	@ [ -f src/binpack.go ] || echo "package main" > src/binpack.go

