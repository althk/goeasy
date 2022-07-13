package grpcutils

import (
	"net/http"
	"os"
	"time"

	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var DefaultKAEP = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var DefaultKASP = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               2 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

var DefaultKACP = keepalive.ClientParameters{
	Time:                8 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             2 * time.Second, // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,            // send pings even without active streams
}

// GRPCServerConfig is used to define the configuration
// for a new GRPC server.
type GRPCServerConfig struct {
	*TLSConfig
	SkipReflection   bool
	SkipHealthServer bool
	SkipZPages       bool
	ZPagesAddr       string
	*KeepAliveConfig
}

type KeepAliveConfig struct {
	KASP          keepalive.ServerParameters
	KACP          keepalive.ClientParameters
	KAEP          keepalive.EnforcementPolicy
	SkipKeepAlive bool
}

func (g *GRPCServerConfig) GetGRPCServerOpts() ([]grpc.ServerOption, error) {
	creds, err := g.TLSConfig.Creds()
	if err != nil {
		return nil, err
	}
	opts := []grpc.ServerOption{
		g.getServerInterceptorChain(),
		grpc.Creds(creds),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	}

	if !g.SkipKeepAlive {
		opts = append(opts,
			grpc.KeepaliveEnforcementPolicy(g.getKAEP()),
			grpc.KeepaliveParams(g.getKASP()),
		)
	}
	return opts, nil
}

func (g *GRPCServerConfig) getKAEP() keepalive.EnforcementPolicy {
	if g.KeepAliveConfig == nil || (keepalive.EnforcementPolicy{}) == g.KeepAliveConfig.KAEP {
		return DefaultKAEP
	}
	return g.KeepAliveConfig.KAEP
}

func (g *GRPCServerConfig) getKACP() keepalive.ClientParameters {
	if g.KeepAliveConfig == nil || (keepalive.ClientParameters{}) == g.KeepAliveConfig.KACP {
		return DefaultKACP
	}
	return g.KeepAliveConfig.KACP
}

func (g *GRPCServerConfig) getKASP() keepalive.ServerParameters {
	if g.KeepAliveConfig == nil || (keepalive.ServerParameters{}) == g.KeepAliveConfig.KASP {
		return DefaultKASP
	}
	return g.KeepAliveConfig.KASP
}

func (g *GRPCServerConfig) GetGRPCDialOpts() ([]grpc.DialOption, error) {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		return nil, err
	}
	creds, err := g.TLSConfig.Creds()
	if err != nil {
		return nil, err
	}
	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(g.getKACP()),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithBlock(),
		getUnaryClientTraceInterceptor(),
	}, nil

}

func (g *GRPCServerConfig) getServerInterceptorChain() grpc.ServerOption {
	logger := zerolog.New(os.Stdout)
	return middleware.WithUnaryServerChain(
		tags.UnaryServerInterceptor(),
		logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger)),
		getUnaryServerTraceInterceptor(),
	)
}

func (g *GRPCServerConfig) NewGRPCServer() (*grpc.Server, error) {
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		return nil, err
	}
	serverOpts, err := g.GetGRPCServerOpts()
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(serverOpts...)

	// register other servers
	if !g.SkipHealthServer {
		h := health.NewServer()
		hpb.RegisterHealthServer(s, h)
		h.Resume()
	}
	if !g.SkipReflection {
		reflection.Register(s)
	}

	if !g.SkipZPages {
		go func() {
			mux := http.NewServeMux()
			zpages.Handle(mux, "/debug")

			if err := http.ListenAndServe(g.ZPagesAddr, mux); err != nil {
				log.Error().Err(err).Msg("Failed to start metrics handler")
			}
		}()
	}
	return s, nil
}
