FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY internal/ ./internal/
RUN CGO_ENABLED=0 go build -o /user-service .

FROM alpine:3.21
COPY --from=build /user-service /user-service
EXPOSE 8080
ENTRYPOINT ["/user-service"]
