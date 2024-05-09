package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"main/internal/domain/models"
	"time"
)

func New(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	// TODO: secret key to yaml
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
