// main.go
package main

import (
	"context"
	"flag"
	"github.com/Rickykn/user-service/db"
	"github.com/Rickykn/user-service/src/handler"
	"github.com/Rickykn/user-service/src/repository"
	"github.com/Rickykn/user-service/src/service"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/Rickykn/drug-proto/gen/user"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// DATABASE CONNECTION
	dbCon := db.Config{
		AppInfo:     "user-service",
		Username:    "postgres",
		Password:    "postgres",
		Database:    "drug-store",
		Host:        "localhost",
		SSLMode:     "disable",
		Port:        5432,
		ConnMaxOpen: 10,
		ConnMaxIdle: 2,
		Logging:     true,
	}

	newDB, err := db.NewDB(dbCon)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userRepo := repository.NewUserRepository(newDB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	user.RegisterUserServiceServer(grpcServer, userHandler)

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	err = user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting HTTP gateway on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
