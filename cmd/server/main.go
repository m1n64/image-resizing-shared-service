package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"image-resizing-shared/internal/delivery/grpc/handlers"
	images "image-resizing-shared/internal/delivery/grpc/pb"
	"image-resizing-shared/internal/delivery/rest"
	"image-resizing-shared/pkg/config"
	"image-resizing-shared/pkg/di"
	"log"
	"net"
	"net/http"
	"os"
)

var dependencies *di.Dependencies
var GinMode string

func init() {
	configPath := resolveConfigPath()

	configuration := config.LoadConfig(configPath)
	dependencies = di.InitDependencies(configuration)
}

func main() {
	go func() {
		serverPort := dependencies.Config.WebServerPort

		fmt.Println("Web server started on port", serverPort, "...")

		if GinMode != "" {
			log.Println("[INFO] Setting Gin mode:", GinMode)
			gin.SetMode(GinMode)
		}

		r := gin.Default()

		rest.RegisterRoutes(r, dependencies)

		if err := r.Run(fmt.Sprintf("0.0.0.0:%s", serverPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	go func() {
		grpcPort := dependencies.Config.GrpcServerPort

		fmt.Println("gRPC server started on port", grpcPort, "...")

		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
		if err != nil {
			log.Fatalf("failed to listen %v", err)
		}

		grpcToken := dependencies.Config.GrpcToken

		var grpcServer *grpc.Server
		if grpcToken != "" {
			grpcServer = grpc.NewServer(
				grpc.UnaryInterceptor(tokenAuthInterceptor(grpcToken)),
			)
		} else {
			grpcServer = grpc.NewServer()
		}

		imagesHandler := handlers.NewImageGRPCHandler(dependencies.ImageService)
		images.RegisterImageServiceServer(grpcServer, imagesHandler)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	select {}
}

func tokenAuthInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader, exists := md["authorization"]
		if !exists || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		if authHeader[0] != token {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

func resolveConfigPath() string {
	args := os.Args
	newArgs := []string{args[0]}
	var configPath string

	for i := 1; i < len(args); i++ {
		if args[i] == "--config" || args[i] == "-c" {
			if i+1 < len(args) {
				configPath = args[i+1]
				i++
				continue
			}
		}
		newArgs = append(newArgs, args[i])
	}

	os.Args = newArgs
	return configPath
}
