# syntax=docker/dockerfile:1

# Build app
FROM golang:1.16-alpine AS build
WORKDIR /app
COPY . .
RUN go build  -o /luno-stream   ./luno-stream/main.go


# Deploy app
FROM gcr.io/distroless/base-debian10
COPY --from=build /luno-stream /luno-stream
USER nonroot:nonroot
ENTRYPOINT ["/luno-stream"]
