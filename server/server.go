package server

import (
	"context"
	"net"
	"strconv"
	"strings"

	"github.com/anant-sharma/go-utils"
	newrelictracing "github.com/anant-sharma/go-utils/new-relic/tracing"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Start - Method to Start gRPC server
func Start(
	config Config,
	protoService func(*grpc.Server),
	opts ...grpc.ServerOption,
) (error, *grpc.Server, net.Listener) {
	hostAddress := strings.Join([]string{
		config.Grpc.Host,
		strconv.Itoa(config.Grpc.Port),
	}, ":")

	lis, err := net.Listen("tcp", hostAddress)
	if err != nil {
		log.Fatal("Unable to listen on", hostAddress)
		return err, nil, nil
	}

	s := grpc.NewServer(
		withServerUnaryInterceptor(),
	)
	protoService(s)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("Starting Server on", hostAddress)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("Unable to start server: ", err)
			return
		}
	}()
	log.Println("Started Server on", hostAddress)

	return nil, s, lis
}

func withServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	txn, c := newrelictracing.NewTransaction(ctx, info.FullMethod)
	txn.AddAttribute("X-Request-ID", utils.GenerateShortID())
	txn.AddAttribute("Method", info.FullMethod)
	defer txn.End()

	// Calls the handler
	h, err := handler(c, req)

	return h, err
}
