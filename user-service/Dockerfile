FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /user-service .

FROM alpine:3.21
COPY --from=build /user-service /user-service
ENTRYPOINT ["/user-service"]
