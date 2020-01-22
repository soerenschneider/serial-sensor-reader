package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"io"
	"strings"
)

const (
	delimiter = '\x0a'
)

type serialReader struct {
	config *serial.Config
	scrape bool
}

func NewSerialReader(device string, baud int) (*serialReader, error) {
	c := &serial.Config{Name: device, Baud: baud}
	if err := testConfig(c); err != nil {
		return nil, err
	}

	return &serialReader{config: c, scrape: true}, nil
}

func testConfig(config *serial.Config) error {
	if config == nil {
		return fmt.Errorf("supplied nil sensorReaderConfig")
	}

	s, err := serial.OpenPort(config)
	if err == nil {
		s.Close()
	}
	return err
}

func (serialReader *serialReader) Interrupt() {
	serialReader.scrape = false
}

func (serialReader *serialReader) ReadSensorData(output chan string) {
	s, err := serial.OpenPort(serialReader.config)
	defer s.Close()
	if err != nil {
		log.Errorf("Received error %s", err.Error())
		return
	}

	err = serialReader.fromReader(output, s)
	if err != nil {
		close(output)
		log.Fatal("Error reading serial data, quitting: %s", err.Error())
	}
}

func (serialReader *serialReader) fromReader(output chan string, rd io.Reader) error {
	for serialReader.scrape {
		reader := bufio.NewReader(rd)
		buf, err := reader.ReadBytes(delimiter)
		if err != nil {
			return err
		}

		read := strings.TrimSpace(string(buf))
		log.Debugf("Read value %s", read)

		output <- read
	}
	
	return nil
}