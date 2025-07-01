package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ufoV1 "github.com/olezhek28/microservices-course-examples/week_4/config/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/api/ufo/v1"
	"github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/config"
	ufoRepository "github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/repository/ufo"
	ufoService "github.com/olezhek28/microservices-course-examples/week_4/config/ufo/internal/service/ufo"
)

const configPath = "./deploy/compose/ufo/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// Создаем MongoDB клиент
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	// Проверяем подключение к MongoDB
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
		return
	}
	log.Println("✅ Connected to MongoDB")

	lis, err := net.Listen("tcp", config.AppConfig().UFOGRPC.Address())
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
	repo := ufoRepository.NewRepository(mongoClient)
	service := ufoService.NewService(repo)
	api := ufoV1API.NewAPI(service)

	ufoV1.RegisterUFOServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().UFOGRPC.Address())
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
