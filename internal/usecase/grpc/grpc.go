package grpc

import (
	"context"
	"fmt"
	school_proto "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/config"
	"github.com/s21platform/school-service/internal/service/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

type SchoolService struct {
	school_proto.UnimplementedSchoolServiceServer
	Lis    net.Listener
	Server *grpc.Server
}

func (s *SchoolService) Login(ctx context.Context, request *school_proto.LoginRequest) (*school_proto.LoginResponse, error) {
	token, err := auth.LoginToPlatform(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &school_proto.LoginResponse{Token: token}, nil
}

func MustSchoolService(cfg *config.Config) *SchoolService {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	s := grpc.NewServer()
	service := &SchoolService{Lis: lis, Server: s}
	school_proto.RegisterSchoolServiceServer(s, service)
	return service
}
