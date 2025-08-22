package main

import (
	"context"
	"fmt"
	"github.com/olezhek28/microservices-course-examples/week_7/easy_metrics/internal/metric"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/olezhek28/microservices-course-examples/week_7/easy_metrics/internal/interceptor"
	noteV1 "github.com/olezhek28/microservices-course-examples/week_7/easy_metrics/pkg/note_v1"
)

const grpcPort = 50051

type server struct {
	noteV1.UnimplementedNoteV1Server
}

// Get ...
func (s *server) Get(ctx context.Context, req *noteV1.GetRequest) (*noteV1.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	// rand.Intn(max - min) + min
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) //nolint:gosec,forbidigo

	return &noteV1.GetResponse{
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
	}, nil
}

func main() {
	ctx := context.Background()

	// Инициализация OpenTelemetry
	meterProvider, err := metric.InitOpenTelemetry()
	if err != nil {
		log.Printf("failed to init OpenTelemetry: %v\n", err.Error())
		return
	}
	defer func() {
		if err = meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter provider: %v\n", err.Error())
		}
	}()

	// Инициализация метрик
	err = metric.Init(ctx)
	if err != nil {
		log.Printf("failed to init metrics: %v\n", err.Error())
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err.Error())
		return
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.MetricsInterceptor,
		),
	)
	reflection.Register(s)
	noteV1.RegisterNoteV1Server(s, &server{})

	log.Printf("server listening at %v\n", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err.Error())
		return
	}
}
