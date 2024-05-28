# Video-Compressor-Tool

This repository contains a small tool wrapper written in go for ffmpeg.
It allows you to compress videos with a simple command and copy their metadata.

I used this tool to compress old videos to minimize the storage space for the cloud.

# Requirements for CLI

- ffmpeg
- exiftool

# Usage with docker

If you don't want to install the cli tools you can use the docker image.
Replace `input.avi` with your filename or directory.

```bash
docker run -v .:/data -it --rm video-compressor-tool:latest -o /data/ /data/input.avi
```

On macOS you can use the hardware acceleration encoder with the following command:

```bash
docker run -v .:/data -it --rm video-compressor-tool:latest -encoder hevc_videotoolbox -o /data/ /data/input.avi
```

# Build

```bash
go build -o video-compressor-tool ./cmd/video-compressor
```