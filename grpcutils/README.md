#

## gprcutils

Provides handy functions that help setting up a new GRPC Server extremely simple.

### Example

The example below sets up a GRPC server with the following options enabled:

* mutual TLS or server side TLS
* zerolog GRPC server interceptor
* client and server keepalive setup
* OpenCensus basic metrics (rpc counts)
* OpenCensus zpages handler for viewing metrics
* OpenTelemetry default tracer with export to opentelemetry collector with support for end-to-end remote traces.

```go
import (
  "google.golang.org/grpc/keepalive"
  "github.com/althk/goeasy/grpcutils"
)

func main() {
  var tlsConfig = &grpcutils.TLSConfig{
    CertFilePath:     "path/to/crt",
    KeyFilePath:      "path/to/key",
    ClientCAFilePath: "path/to/clientca.crt",
    RootCAFilePath:   "path/to/rootca.crt",
    SkipTLS:          false, // Setting this to true returns grpc.Server with Insecure creds
    NoClientCert:     false, // Setting this to true turns off mutual TLS auth and does only server auth
  }

  var kaConfig = &grpcutils.KeepAliveConfig{
    KASP:          keepalive.ServerParameters{},  // Skipping this or sending the empty struct will initialize with default values
    KACP:          keepalive.ClientParameters{}, // Skipping this or sending the empty struct will initialize with default values
    KAEP:          keepalive.EnforcementPolicy{}, // Skipping this or sending the empty struct will initialize with default values
    SkipKeepAlive: false,  // Setting this to true skips KeepAlive
  }

  var grpcConfig = &grpcutils.GRPCServerConfig{
    TLSConfig:        tlsConfig,
    SkipReflection:   false,
    SkipHealthServer: false,
    SkipZPages:       false,
    ZPagesAddr:       "localhost:5555",
    KeepAliveConfig: kaConfig,
  }

  grpcServer, err := grpcConfig.NewGRPCServer()
  if err != nil {
    // handle error
  }
  // register your service with grpcServer
  pb.RegisterMyServiceServer(grpcServer, *myServiceServer)

  // start an OpenTelemetry Tracer
  // this will export the traces to an OpenTelemetry Collector server running at
  // otelcollector1:4317
  tp, err := grpcutils.OTelTraceProvider("my-service-name", "otelcollector1:4317")
  if err != nil {
    // handle error
  }
  defer func() {
    if err := tp.Shutdown(context.Background()); err != nil {
     log.Printf("Error shutting down tracer provider: %v", err)
    }
  }()
}
```
