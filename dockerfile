# Start from a base image with Go installed
FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /main ./cmd/thumbnail-server/main.go

EXPOSE 8080

CMD ["./main"]