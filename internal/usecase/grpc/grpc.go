package grpc

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
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

func (s *SchoolService) Login(ctx context.Context, request *school_proto.SchoolLoginRequest) (*school_proto.SchoolLoginResponse, error) {
	// Setup client. Switch off redirects
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	token, err := auth.LoginToPlatform(client, "https://auth.sberclass.ru", request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &school_proto.SchoolLoginResponse{Token: token}, nil
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
