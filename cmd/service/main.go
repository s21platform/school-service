package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	logger_lib "github.com/s21platform/logger-lib"
	school_proto "github.com/s21platform/school-proto/school-proto"
	"github.com/s21platform/school-service/internal/config"
	"github.com/s21platform/school-service/internal/infra"
	"github.com/s21platform/school-service/internal/repository/redis"
	"github.com/s21platform/school-service/internal/service"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	redisRepo := redis.New(cfg)
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(infra.Logger(logger)),
	)

	srv := service.New(redisRepo)
	school_proto.RegisterSchoolServiceServer(s, srv)

	log.Println("starting server, port:", cfg.Service.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("Cannnot listen port: %s; Error: %s", cfg.Service.Port, err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Cannnot start service: %s; Error: %s", cfg.Service.Port, err)
	}
}
