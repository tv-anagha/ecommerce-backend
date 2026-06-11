FROM golang:1.25-alpine AS build

ENV GOMAXPROCS=1
ENV GOGC=50

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -p 1 -o /user-service ./cmd/server

FROM alpine:3.21
COPY --from=build /user-service /user-service
EXPOSE 8082
ENTRYPOINT ["/user-service"]
