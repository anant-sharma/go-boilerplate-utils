package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func StopWatch(grpcServer *grpc.Server, lis net.Listener, gwServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Stop Server
	log.Println("Shutting Down Server ...")
	grpcServer.Stop()

	// Close Listener
	log.Println("Closing Listener ...")
	lis.Close()

	// Close Gateway Server
	log.Println("Closing Gateway ...")
	gwServer.Shutdown(context.Background())

	log.Println("Server Shutdown Complete.")
}
