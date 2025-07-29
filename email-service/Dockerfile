FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o email-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/email-service /app/email-service
COPY .env .env

EXPOSE 8084
CMD ["./email-service"]
