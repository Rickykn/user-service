// main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Rickykn/user-service/db"
	"github.com/Rickykn/user-service/src/config"
	"github.com/Rickykn/user-service/src/handler"
	"github.com/Rickykn/user-service/src/interceptor"
	"github.com/Rickykn/user-service/src/logger"
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
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("grpc: main failed to load and parse config: %s", err)
		return
	}

	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// DATABASE CONNECTION
	dbCon := db.Config{
		AppInfo:     cfg.AppInfo.Name,
		Username:    cfg.PostgreSQL.Username,
		Password:    cfg.PostgreSQL.Password,
		Database:    cfg.PostgreSQL.Database,
		Host:        cfg.PostgreSQL.Host,
		SSLMode:     cfg.PostgreSQL.SSLMode,
		Port:        cfg.PostgreSQL.Port,
		ConnMaxOpen: cfg.PostgreSQL.MaxOpenConns,
		ConnMaxIdle: cfg.PostgreSQL.MaxIdleConns,
		Logging:     cfg.PostgreSQL.Logging,
	}

	newDB, err := db.NewDB(dbCon)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	logger.Init(cfg.AppInfo.Name, cfg.AppInfo.Environment, cfg.Logger.Level)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor()),
	)
	addr := fmt.Sprintf(":%d", cfg.GRPCServer.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	userRepo := repository.NewUserRepository(newDB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)
	user.RegisterUserServiceServer(grpcServer, userHandler)

	go func() {
		log.Printf("Starting gRPC server on :%d", cfg.GRPCServer.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	err = user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	addrRestProxy := fmt.Sprintf(":%d", cfg.GRPCRESTProxyServer.Port)
	log.Printf("Starting HTTP gateway on %s", addrRestProxy)
	if err := http.ListenAndServe(addrRestProxy, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
