.PHONY: build
build:
	if [ ! -d build ]; then mkdir build; fi
	go build -o build/serial-sensor-reader
