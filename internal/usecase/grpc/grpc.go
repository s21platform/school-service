package grpc

//
//import (
//	"context"
//	"github.com/go-resty/resty/v2"
//	school_proto "github.com/s21platform/school-proto/school-proto"
//	"github.com/s21platform/school-service/internal/service"
//	"google.golang.org/grpc"
//	"net"
//)
//
//type SchoolService struct {
//	school_proto.UnimplementedSchoolServiceServer
//	Lis    net.Listener
//	Server *grpc.Server
//}
//
//func (s *SchoolService) Login(ctx context.Context, request *school_proto.SchoolLoginRequest) (*school_proto.SchoolLoginResponse, error) {
//	// Setup client. Switch off redirects
//	client := resty.New()
//	client.SetRedirectPolicy(resty.NoRedirectPolicy())
//
//	token, err := service.LoginToPlatform(client, "https://auth.sberclass.ru", request.Email, request.Password)
//	if err != nil {
//		return nil, err
//	}
//	return &school_proto.SchoolLoginResponse{Token: token}, nil
//}
