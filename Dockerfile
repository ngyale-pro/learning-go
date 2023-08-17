# Build stage
FROM golang:1.20-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration


# Not really opening this port, rather documenting it
EXPOSE 8080
# When using CMD + Entrypoint together, the CMD is given to the entrypoint as a parameter ENTRYPOINT [ "/app/start.sh", "/app/main"]
CMD ["/app/main"] 
ENTRYPOINT [ "/app/start.sh" ]

