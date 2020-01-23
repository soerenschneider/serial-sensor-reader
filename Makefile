.PHONY: build
build:
	if [ ! -d build ]; then mkdir build; fi
	CGO_ENABLED=0 go build -o build/serial-sensor-reader ./serialsensorreader 

raspberry:
	if [ ! -d build ]; then mkdir build; fi
	env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o build/serial-sensor-reader.arm7 ./serialsensorreader 
	env GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=0 go build -o build/serial-sensor-reader.arm5 ./serialsensorreader 
