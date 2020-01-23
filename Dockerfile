FROM golang:1.13 as builder

COPY . /vgo/
WORKDIR /vgo/
RUN ls 
RUN make build

FROM alpine:3 as runtime
COPY --from=builder /vgo/build/serial-sensor-reader /serial-sensor-reader
USER 65534

ENTRYPOINT ["/serial-sensor-reader"]
