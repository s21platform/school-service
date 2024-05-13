package main

import (
	"fmt"
	"github.com/s21platform/school-service/internal/config"
	"github.com/s21platform/school-service/internal/usecase/grpc"
	"log"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	service := grpc.MustSchoolService(cfg)
	fmt.Println("Starting service")
	if err := service.Server.Serve(service.Lis); err != nil {
		log.Fatalf("Error while starting service: %s", err)
	}
}
