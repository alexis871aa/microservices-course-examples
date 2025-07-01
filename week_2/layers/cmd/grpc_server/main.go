package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ufoV1API "github.com/olezhek28/microservices-course-examples/week_2/layers/internal/api/ufo/v1"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_2/layers/internal/repository/ufo"
	ufoService "github.com/olezhek28/microservices-course-examples/week_2/layers/internal/service/ufo"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/layers/pkg/proto/ufo/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	repo := ufoRepository.NewRepository()
	service := ufoService.NewService(repo)
	api := ufoV1API.NewAPI(service)

	ufoV1.RegisterUFOServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
