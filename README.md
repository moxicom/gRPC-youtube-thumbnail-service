# gRPC-youtube-thumbnail-service

## Overview
The gRPC YouTube Thumbnail Service is a microservice that provides a gRPC API to fetch thumbnails for YouTube videos. 
This service accepts a YouTube video URLs and returns the corresponding thumbnail URL. 
It's built using gRPC and can be integrated into larger systems that require thumbnail retrieval.

## Features
Fetches YouTube video thumbnails using video URLs.
Provides a gRPC API for easy integration with other services.
Error handling for invalid or non-existent video URLs.

Using runtime cache

## Setup server
To setup server use your local machine with go or DOCKER COMPOSE.
command
``` bash
go run cmd/thumbnail-server/main.go

# by docker
docker-compose up
```
Flags: 
- `envLog` values `'local'` or `'prod'`

## Using client
To run client use flags
- `-server_addr`: The server address in the format `host:port`. Default: `localhost:8080`.
- `-async`: Asynchronous mode for downloading thumbnails. Default: `false`.
- `-envLog`: Sets the logging environment. Possible values: `'prod'` or `'local'`. Default: `'prod'`.

```bash
go run cmd/thumbnail-client/main.go [flags] [arguments] 

# example
go run cmd/thumbnail-client/main.go \
-server_addr="localhost:8080" \
-async=true \
-envLog="prod" \
https://www.youtube.com/watch?v=dQw4w9WgXcQ \
https://www.youtube.com/watch?v=dQw4w9WgXcQ
```
## Path to save
Client saves thumbnails to `./output` folder.

## Setting up CI/CD Pipeline with GitHub Actions
This project uses GitHub Actions to automate the build, test, and deployment process of a Go application using Docker and Docker Compose. This section explains how to set up and use the pipeline.

## Configure Secrets in GitHub
- `DOCKER_REGISTRY`: The address of your Docker Registry (e.g., docker.io for Docker Hub).
- `SSH_PRIVATE_KEY`: The private SSH key for accessing your server.
- `SERVER_USER`: The username on your server.
- `SERVER_HOST`: The address of your server.

