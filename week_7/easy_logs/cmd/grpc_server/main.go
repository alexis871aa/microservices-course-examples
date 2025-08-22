package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	applogger "github.com/olezhek28/microservices-course-examples/week_7/easy_logs/internal/logger"
	noteV1 "github.com/olezhek28/microservices-course-examples/week_7/easy_logs/pkg/proto/note_v1"
)

const grpcPort = 50051

type server struct {
	noteV1.UnimplementedNoteV1Server
}

// Get ...
func (s *server) Get(ctx context.Context, req *noteV1.GetRequest) (*noteV1.GetResponse, error) {
	applogger.Info(ctx, "Get request", zap.Int64("id", req.GetId()))
	if req.GetId() == 0 {
		applogger.Error(ctx, "invalid id", zap.String("error", "id is empty"))
		return nil, errors.Errorf("id is empty")
	}

	// rand.Intn(max - min) + min
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) //nolint:gosec,forbidigo

	resp := &noteV1.GetResponse{
		Note: &noteV1.Note{
			Id: req.GetId(),
			Info: &noteV1.NoteInfo{
				Title:    gofakeit.BeerName(),
				Context:  gofakeit.IPv4Address(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}
	applogger.Info(
		ctx,
		"Get response",
		zap.Int64("id", resp.GetNote().GetId()),
		zap.String("title", resp.GetNote().GetInfo().GetTitle()),
		zap.String("context", resp.GetNote().GetInfo().GetContext()),
		zap.String("author", resp.GetNote().GetInfo().GetAuthor()),
		zap.Bool("is_public", resp.GetNote().GetInfo().GetIsPublic()),
	)

	return resp, nil
}

func main() {
	ctx := context.Background()

	// Инициализируем логгер с Tee архитектурой (stdout + OTLP)
	// enableOTLP=true включает автоматическую отправку в OpenTelemetry коллектор
	if err := applogger.Init("info", true, true); err != nil {
		panic(err)
	}
	// Обеспечиваем корректное завершение работы логгера
	defer func() {
		_ = applogger.Sync()  // Сбрасываем буферы zap
		_ = applogger.Close() // Закрываем OTLP ресурсы
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		applogger.Error(ctx, "failed to listen", zap.String("error", err.Error()))
		return
	}

	s := grpc.NewServer()
	reflection.Register(s)
	noteV1.RegisterNoteV1Server(s, &server{})

	applogger.Info(ctx, "server listening", zap.String("addr", lis.Addr().String()))

	go func() {
		err = s.Serve(lis)
		if err != nil {
			applogger.Error(ctx, "failed to serve", zap.String("error", err.Error()))
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(quit)

	applogger.Info(ctx, "🛑 Shutting down gRPC server...")
	s.GracefulStop()
	applogger.Info(ctx, "✅ Server stopped")
}
