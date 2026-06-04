FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o /product-service ./cmd/server

FROM alpine:3.21
COPY --from=build /product-service /product-service
EXPOSE 8081
ENTRYPOINT ["/product-service"]
