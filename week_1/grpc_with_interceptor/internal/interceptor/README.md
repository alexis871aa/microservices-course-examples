# gRPC интерцепторы

## Общее описание

В этом пакете реализованы gRPC интерцепторы, которые можно использовать для добавления дополнительной функциональности к gRPC серверу.

## LoggerInterceptor

`LoggerInterceptor` - это унарный серверный интерцептор, который логирует информацию о начале и окончании вызовов методов gRPC, а также измеряет время выполнения каждого метода.

### Особенности

- Логирует начало вызова метода с его именем
- Измеряет время выполнения метода
- Логирует успешное завершение метода с временем выполнения
- В случае ошибки логирует код ошибки, сообщение и время выполнения

### Пример использования

```go
import (
    "github.com/olezhek28/microservices-course-examples/week_1/grpc_with_interceptor/internal/interceptor"
    "google.golang.org/grpc"
)

// Создание gRPC сервера с логирующим интерцептором
server := grpc.NewServer(
    grpc.UnaryInterceptor(interceptor.LoggerInterceptor()),
)
```

## Использование нескольких интерцепторов

Для использования нескольких интерцепторов одновременно можно использовать пакеты, такие как `github.com/grpc-ecosystem/go-grpc-middleware`:

```go
import (
    "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/olezhek28/microservices-course-examples/week_1/grpc_with_interceptor/internal/interceptor"
    "google.golang.org/grpc"
)

// Создание gRPC сервера с несколькими интерцепторами
server := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        interceptor.LoggerInterceptor(),
        interceptor.LogErrorInterceptor(),
        // Другие интерцепторы...
    )),
)
``` 