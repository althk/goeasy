// Package grpcutils implements utility functions for simplying GRPC setup boilerplate.
//
// * Simplified GRPC Server creation including:
//   * TLS setup
//   * Interceptors
//   * Opentelemetry (basic metrics and tracing interceptors)
//
// Example:
// The code snippet below does the following setup:
// * TLS Setup (insecure, server auth, mutual auth)
// * grpc logging interceptor (zerolog, which is the most performant of all the logging libs)
// * grpc keep alive config
// * OpenTelemetry tracing via client and server interceptors and exporting to OpenTelemetry Collector
// * OpenCensus basic metrics, exposed over the ZPagesAddr
//
// import (
// 	"google.golang.org/grpc/keepalive"
// 	"github.com/althk/goeasy/grpcutils"
// )

// func main() {
// 	var tlsConfig = &TLSConfig{
// 		CertFilePath:     "path/to/crt",
// 		KeyFilePath:      "path/to/key",
// 		ClientCAFilePath: "path/to/clientca.crt",
// 		RootCAFilePath:   "path/to/rootca.crt",
// 		SkipTLS:          false, // Setting this to true returns grpc.Server with Insecure creds
// 		NoClientCert:     false, // Setting this to true turns off mutual TLS auth and does only server auth
// 	}

// 	var kaConfig = &KeepAliveConfig{
// 		KASP:          keepalive.ServerParameters{},  // Skipping this or sending the empty struct will initialize with default values
// 		KACP:          keepalive.ClientParameters{}, // Skipping this or sending the empty struct will initialize with default values
// 		KAEP:          keepalive.EnforcementPolicy{}, // Skipping this or sending the empty struct will initialize with default values
// 		SkipKeepAlive: false,  // Setting this to true skips KeepAlive
// 	}

// 	var grpcConfig = &GRPCServerConfig{
// 		TLSConfig:        tlsConfig,
// 		SkipReflection:   false,
// 		SkipHealthServer: false,
// 		SkipZPages:       false,
// 		ZPagesAddr:       "localhost:5555",
// 		KeepAliveConfig: kaConfig,
// 	}

// 	grpcServer, err := grpcConfig.NewGRPCServer()
// 	if err != nil {
// 		// handle error
// 	}
// 	// register your service with grpcServer
// 	pb.RegisterMyServiceServer(grpcServer, *myServiceServer)
// 	// start an OpenTelemetry Tracer
// 	tp, err := grpcutils.OTelTraceProvider("my-service-name", "otelcollector1:4317")
// 	if err != nil {
// 		// handle error
// 	}
// 	defer func() {
// 		if err := shutdownFn(context.Background()); err != nil {
// 			log.Printf("Error shutting down tracer provider: %v", err)
// 		}
// 	}()
// }
package grpcutils
