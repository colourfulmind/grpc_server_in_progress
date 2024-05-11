package grpcclient

import (
	"context"
	"github.com/gin-gonic/gin"
	"main/internal/domain/models"
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
	ginCtx.IndentedJSON(http.StatusCreated, user)

	resp, err := c.sso.RegisterNewUser(context.Background(), &server.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		return
	}

	c.log.Info("user is created")
	ginCtx.IndentedJSON(http.StatusCreated, resp.UserId)
}

func (c *Client) GetUser(ginCtx *gin.Context) {
	var user models.NewUser

	// TODO: return error
	if err := ginCtx.BindJSON(&user); err != nil {
		return
	}

	resp, err := c.sso.Login(context.Background(), &server.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		ginCtx.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.log.Info("user is gotten")
	ginCtx.IndentedJSON(http.StatusOK, resp)
}
