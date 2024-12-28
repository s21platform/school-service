package service

import (
	"context"
	"errors"
	"fmt"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/school-service/internal/config"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	school "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/usecase/edu_school"
)

type Server struct {
	school.UnimplementedSchoolServiceServer
	redisR RedisR
}

func New(redis RedisR) *Server {
	return &Server{redisR: redis}
}

func (s *Server) Login(ctx context.Context, request *school.SchoolLoginRequest) (*school.SchoolLoginResponse, error) {
	log.Println("Try to get school token for: ", request.Email)
	resp, err := edu_school.LoginToPlatform(request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	err = s.redisR.Set(ctx, resp.AccessToken, resp.AccessToken, 10*time.Hour)
	if err != nil {
		log.Println(err)
	}
	return &school.SchoolLoginResponse{Token: resp.AccessToken}, nil
}

func (s *Server) GetCampuses(ctx context.Context, _ *school.Empty) (*school.CampusesOut, error) {
	token, err := s.redisR.Get(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := edu_school.GetAllCampuses(token)
	if err != nil {
		log.Printf("error of getting campuses: %v", err)
		return nil, err
	}

	var needCampuses []*school.Campus
	for _, value := range resp.Campuses {
		needCampuses = append(needCampuses, &school.Campus{
			CampusUuid: value.Uuid,
			ShortName:  value.ShortName,
			FullName:   value.FullName,
		})
	}

	return &school.CampusesOut{Campuses: needCampuses}, nil
}

func (s *Server) GetPeers(ctx context.Context, in *school.GetPeersIn) (*school.GetPeersOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetPeers")

	token, err := s.redisR.Get(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get peer token, err: %v", err))
		return nil, status.Errorf(codes.Internal, "cannot get peer token, err: %v", err)
	}

	peers, err := edu_school.GetPeers(ctx, token, in.CampusUuid, in.Offset, in.Limit)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get peers from edu, err: %v", err))
		return nil, status.Errorf(codes.Internal, "cannot get peers from edu, err: %v", err)
	}

	return &school.GetPeersOut{Peer: peers}, nil
}

func (s *Server) GetTribesByCampusUuid(ctx context.Context, in *school.CampusUuidIn) (*school.TribesOut, error) {
	if in.CampusUuid == "" {
		return nil, errors.New("campus uuid is empty")
	}

	log.Println("Trying to get list of tribes by campus uuid: ", in.CampusUuid)

	token, err := s.redisR.Get(ctx)

	if err != nil {
		return nil, err
	}
	resp, err := edu_school.GetTribesOfCampus(in.CampusUuid, token)
	if err != nil {
		log.Printf("error of getting tribes: %v", err)
		return nil, err
	}

	var needTribes []*school.Tribe
	for _, value := range resp.Tribes {
		needTribes = append(needTribes, &school.Tribe{
			Id:   value.Id,
			Name: value.Name,
		})
	}

	return &school.TribesOut{Tribes: needTribes}, nil
}
