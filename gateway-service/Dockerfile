FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o gateway-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/gateway-service /app/gateway-service
COPY .env .env

EXPOSE 8080
CMD ["./gateway-service"]
