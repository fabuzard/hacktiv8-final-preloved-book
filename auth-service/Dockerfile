FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o auth-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/auth-service /app/auth-service
COPY .env .env

EXPOSE 8080
CMD ["./auth-service"]
