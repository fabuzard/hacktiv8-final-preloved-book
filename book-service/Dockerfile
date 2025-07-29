FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o book-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/book-service /app/book-service
COPY .env .env

EXPOSE 8081
CMD ["./book-service"]
