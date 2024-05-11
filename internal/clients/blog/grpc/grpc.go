package grpcclient

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"main/internal/config"
	server "main/protos/gen/go/blog"
	"time"
)

// Client represents data about a client
type Client struct {
	sso    server.SSOClient
	log    *slog.Logger
	Router *gin.Engine
	config *config.Config
	token  string
}

// New returns a Client struct
func New(cc *grpc.ClientConn, log *slog.Logger, config *config.Config) *Client {
	c := &Client{
		sso:    server.NewSSOClient(cc),
		log:    log,
		Router: gin.Default(),
		config: config,
	}
	//c.ConfigureRouter()
	return c
}

// ConfigureRouter handles requests
func (c *Client) ConfigureRouter() {
	c.Router.POST("/", c.RegisterNewUser)
	c.Router.GET("/", c.GetUser)
}

//func (c *Client) Start() error {
//	return http.ListenAndServe(":8080", c.Router)
//}

func NewConnection(ctx context.Context, log *slog.Logger, addr string, retriesCount int, timeout time.Duration) (*grpc.ClientConn, error) {
	const op = "internal/clients/blog/NewConnection"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cc, nil
}

func InterceptorLogger(log *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), msg, fields...)
	})
}
