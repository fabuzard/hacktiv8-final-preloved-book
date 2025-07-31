# Stage 1: Build
FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o transaction-service main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/transaction-service /app/transaction-service
COPY .env .env

EXPOSE 8083

CMD ["./transaction-service"]
