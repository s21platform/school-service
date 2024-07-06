package main

import (
	"fmt"
	school_proto "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/config"
	"github.com/s21platform/school-service/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.MustLoad()

	srv := service.New()
	s := grpc.NewServer()
	school_proto.RegisterSchoolServiceServer(s, srv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
