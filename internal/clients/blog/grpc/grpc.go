package grpcclent

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log/slog"
	"main/internal/config"
	server "main/protos/gen/go/blog"
)

// Client represents data about a client
type Client struct {
	sso    server.SSOClient
	log    *slog.Logger
	router *gin.Engine
	config *config.Config
	token  string
}

// New returns a Client struct
func New(cc *grpc.ClientConn, log *slog.Logger, config *config.Config) *Client {
	c := &Client{
		sso:    server.NewSSOClient(cc),
		log:    log,
		router: gin.Default(),
		config: config,
	}
	c.ConfigureRouter()
	return c
}

// ConfigureRouter handles requests
func (c *Client) ConfigureRouter() {
	c.router.GET("/")
}
