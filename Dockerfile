FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o gosubscriber ./main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/gosubscriber .

EXPOSE 8080

CMD ["./gosubscriber"]
