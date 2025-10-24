# gRPC

Генерация из `.proto` файла
```shell
protoc --go_out=./ --go-grpc_out=./ api/proto/hello.proto
```

# Golang
Удаление всех кешированных зависимостей
```shell
go clean -modcache
```
## Дерево зависимостей
Установить сторонний инструмент
```shell
go install github.com/PaulXu-cn/go-mod-graph-chart/gmchart@latest
```
Вывод дерева
```shell
go mod graph | gmchart
```

# Stack
- fyne.io
- gRPC
- zap