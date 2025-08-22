package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/google/uuid"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonV1 "github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared/pkg/proto/common/v1"
	iamV1 "github.com/olezhek28/microservices-course-examples/week_8/grpc_with_iam/shared/pkg/proto/iam/v1"
)

const (
	grpcAddr = "0.0.0.0:50052"

	SessionCookieName = "X-Session-Uuid"

	HeaderUserUUID    = "X-User-Uuid"
	HeaderUserLogin   = "X-User-Login"
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "X-Auth-Status"

	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"

	ContentTypeJSON = "application/json"

	AuthStatusDenied = "denied"
)

type IAMServer struct {
	iamV1.UnimplementedIAMServiceServer
	authv3.UnimplementedAuthorizationServer
}

func NewIAMServer() *IAMServer {
	return &IAMServer{}
}

func (s *IAMServer) Whoami(_ context.Context, req *iamV1.WhoamiRequest) (*iamV1.WhoamiResponse, error) {
	log.Printf("Whoami called with session_uuid: %s", req.SessionUuid)

	_, err := uuid.Parse(req.SessionUuid)
	if err != nil {
		log.Printf("Invalid session UUID: %v", err)
		return nil, fmt.Errorf("invalid session UUID: %w", err)
	}

	return &iamV1.WhoamiResponse{
		Session: &commonV1.Session{
			Uuid:      req.SessionUuid,
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
			ExpiresAt: timestamppb.New(gofakeit.Date()),
		},
		User: &commonV1.User{
			Uuid: gofakeit.UUID(),
			Info: &commonV1.UserInfo{
				Login:               gofakeit.Username(),
				Email:               gofakeit.Email(),
				NotificationMethods: []commonV1.NotificationMethod{commonV1.NotificationMethod_NOTIFICATION_METHOD_PUSH},
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *IAMServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	log.Printf("External Authorization Check called")

	sessionUUID, err := s.extractSessionUUID(req)
	if err != nil {
		log.Printf("Session extraction failed: %v", err)
		return s.denyRequest("Missing or invalid session", 403), nil
	}

	log.Printf("Extracted session_uuid: %s", sessionUUID)

	whoamiResp, err := s.Whoami(ctx, &iamV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		log.Printf("Whoami failed: %v", err)
		return s.denyRequest("Invalid session", 403), nil
	}

	return s.allowRequest(whoamiResp), nil
}

func (s *IAMServer) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if cookieHeader, ok := headers[HeaderCookie]; ok && cookieHeader != "" {
		sessionUUID := s.extractSessionFromCookies(cookieHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", fmt.Errorf("session uuid not found in cookies")
}

func (s *IAMServer) extractSessionFromCookies(cookieHeader string) string {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add(HeaderCookie, cookieHeader)

	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		var sessionUUID string
		sessionUUID, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}

		return sessionUUID
	}

	return ""
}

func (s *IAMServer) allowRequest(whoamiResp *iamV1.WhoamiResponse) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: whoamiResp.User.Uuid,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: whoamiResp.User.Info.Login,
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

func (s *IAMServer) denyRequest(message string, statusCode int32) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: typev3.StatusCode(statusCode),
				},
				Body: fmt.Sprintf(`{"error": "%s", "timestamp": "%s"}`,
					message, time.Now().Format(time.RFC3339)),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderAuthStatus,
							Value: AuthStatusDenied,
						},
					},
				},
			},
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err.Error())
		return
	}

	s := grpc.NewServer()

	iamServer := NewIAMServer()
	iamV1.RegisterIAMServiceServer(s, iamServer)
	authv3.RegisterAuthorizationServer(s, iamServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("iam.v1.IAMService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 IAM gRPC server listening on %s\n", grpcAddr)
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
	log.Println("🛑 Shutting down IAM gRPC server...")
	s.GracefulStop()
	log.Println("✅ IAM Server stopped")
}
