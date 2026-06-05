FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o /order-service ./cmd/server

FROM alpine:3.21
COPY --from=build /order-service /order-service
EXPOSE 8084
ENTRYPOINT ["/order-service"]
