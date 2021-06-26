package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func gatewayHandler(mux, oaHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger") {

			if r.URL.Path == "/swagger" {
				http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
				return
			}

			oaHandler.ServeHTTP(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})
}

func StartHTTPProxy(
	config Config,
	openAPIDir string,
	protoHandlerFromEndpoint func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error,
) *http.Server {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	endpoint := strings.Join([]string{
		config.Grpc.Host,
		strconv.Itoa(config.Grpc.Port),
	}, ":")

	// Register Service
	err := protoHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		log.Panicln("Unable to register protoHandlerFromEndpoint")
		return nil
	}

	// Open API Handler
	oaHandler := http.StripPrefix("/swagger/", getOpenAPIHandler(openAPIDir))

	hostAddress := strings.Join([]string{
		config.HTTP.Host,
		strconv.Itoa(config.HTTP.Port),
	}, ":")

	log.Println("Starting HTTP Proxy on", hostAddress)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	gwServer := &http.Server{
		Addr:    hostAddress,
		Handler: apm(tracing(logging(gatewayHandler(mux, oaHandler)))),
	}

	log.Printf("Serving gRPC-Gateway on http://%s", hostAddress)
	log.Printf("Serving OpenAPI Documentation on http://%s/swagger/", hostAddress)

	gwErr := gwServer.ListenAndServe()
	if gwErr != nil {
		log.Panicln("Unable to start HTTP Proxy")
		return nil
	}

	return gwServer
}
