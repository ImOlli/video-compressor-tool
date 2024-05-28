FROM golang:1.22.3-bookworm as build
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o /app/video-compressor /app/cmd/video-compressor

FROM alpine:latest
RUN apk add exiftool ffmpeg
COPY --from=build /app/video-compressor /app/video-compressor
WORKDIR /app

ENTRYPOINT ["/app/video-compressor"]