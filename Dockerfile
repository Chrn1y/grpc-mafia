FROM golang:latest

WORKDIR /grpc-mafia

COPY . .

RUN go mod download

EXPOSE 8000

CMD ["go", "run", "server/main.go"]