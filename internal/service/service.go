package service

import (
	"context"
	"errors"
	school "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/usecase/edu_school"
	"log"
	"time"
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
