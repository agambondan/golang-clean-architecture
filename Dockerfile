FROM golang:latest

LABEL maintainer="Quique agam.pro234@gmail.com"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 5050

RUN go build -o build/main

CMD ["./build/main"]