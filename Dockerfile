# http://www.inanzzz.com/index.php/post/1sfg/multi-stage-docker-build-for-a-golang-application-with-and-without-vendor-directory
ARG debianVersion=latest
FROM debian:${debianVersion} AS rtl433-builder

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    libusb-1.0-0-dev \
    libsoapysdr-dev \
    librtlsdr-dev \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /build

RUN git clone https://github.com/switchdoclabs/rtl_433
WORKDIR ./rtl_433
ARG rtl433GitRevision=master
RUN git checkout ${rtl433GitRevision}
WORKDIR ./build
RUN cmake ..
RUN make -j 4
RUN cat Makefile
WORKDIR /build/root
WORKDIR /build/rtl_433/build
RUN make DESTDIR=/build/root/ install
RUN ls -lah /build/root

# Compile stage
FROM golang:1.16.5 AS go-builder
ENV CGO_ENABLED 0

WORKDIR /weather-sensor-bridge

ADD . ./

RUN make build

# Final stage
ARG debianVersion=latest
FROM debian:${debianVersion} AS output
LABEL org.opencontainers.image.source https://github.com/geoff-coppertop/weather-sensor-bridge

RUN apt-get update && apt-get install -y \
    libusb-1.0-0 \
    librtlsdr-dev \
    libsoapysdr-dev \
    soapysdr-module-all \
 && rm -rf /var/lib/apt/lists/*

COPY --from=rtl433-builder /build/root/ /
COPY --from=go-builder /weather-sensor-bridge/bin/weather-sensor-bridge /

# Run
CMD ["/weather-sensor-bridge"]
