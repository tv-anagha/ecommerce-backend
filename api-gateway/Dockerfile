FROM golang:1.25-alpine AS build

ENV GOMAXPROCS=1
ENV GOGC=50

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -p 1 -o /api-gateway ./cmd/server

FROM alpine:3.21
COPY --from=build /api-gateway /api-gateway
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
