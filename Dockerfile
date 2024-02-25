FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

RUN GOOS=linux GOARCH=amd64 go build -o ./gosubscriber ./main.go


CMD ["go", "run", "."]