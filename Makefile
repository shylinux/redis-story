all:
	@echo && date
	go build -o ice.bin main.go && chmod u+x ice.bin && ./ice.sh restart
