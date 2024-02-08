# Use an official Node.js runtime as the base image
FROM golang:1.22-alpine as builder

WORKDIR /tmp/build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/build/proxy-server ./cmd/ffmpeg-proxy/main.go

FROM alpine:3.19

RUN apk update && apk add ffmpeg linux-headers alsaconf alsa-lib alsa-utils

RUN echo "snd_pcm" >> /etc/modules && \
    echo "snd_seq" >> /etc/modules && \
    echo "snd_seq_device" >> /etc/modules && \
    echo "snd_usb_audio" >> /etc/modules

# Set the working directory in the container to /app
WORKDIR /app

# Copy the rest of the application's code to the working directory
COPY --from=builder /tmp/build/proxy-server /app/server

# Make port 8081 available outside the container
EXPOSE 8080

# Start the application when the container is run
CMD [ "./server" ]
