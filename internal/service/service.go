package service

import (
	"context"
	school_proto "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/usecase/edu_school"
	"log"
)

type Server struct {
	school_proto.UnimplementedSchoolServiceServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) Login(ctx context.Context, request *school_proto.SchoolLoginRequest) (*school_proto.SchoolLoginResponse, error) {
	log.Println("Try to get school token for: ", request.Email)
	resp, err := edu_school.LoginToPlatform(request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &school_proto.SchoolLoginResponse{Token: resp.AccessToken}, nil
}
