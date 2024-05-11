package grpcclient

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/internal/domain/models"
	"main/pkg/logger/sl"
	server "main/protos/gen/go/blog"
	"net/http"
)

// RegisterNewUser registers new user
func (c *Client) RegisterNewUser(ginCtx *gin.Context) {
	var user models.NewUser

	// TODO: return error
	if err := ginCtx.BindJSON(&user); err != nil {
		return
	}

	c.log.Info("user", user)

	resp, err := c.sso.RegisterNewUser(context.Background(), &server.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	// "error": "rpc error: code = Unavailable desc = connection error: desc = \"error reading server preface: http2: frame too large\""
	if err != nil {
		c.log.Error("user was not created", sl.Err(err))
		return
	}
	ginCtx.IndentedJSON(http.StatusCreated, resp.UserId)
}

func (c *Client) GetUser(ginCtx *gin.Context) {
	var user = models.NewUser{Email: "test@test.test", Password: "1234"}

	resp, err := c.sso.Login(context.Background(), &server.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		ginCtx.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	ginCtx.IndentedJSON(http.StatusOK, resp)
}
