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