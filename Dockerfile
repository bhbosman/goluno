FROM golang
WORKDIR /app
COPY . .
RUN go build ./luno-stream/main.go
