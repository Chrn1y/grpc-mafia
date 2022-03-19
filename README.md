# grpc-mafia

Запуск клиента:
```shell
go run client/main.go
```
Генерация proto:
```shell
protoc -I ./proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative ./proto/mafia.proto
```