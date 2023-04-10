package auth

import (
	structures "bcraft/api/structures"
	requests "bcraft/api/structures/requests"
	"bcraft/db"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	salt      = "nvasjirbfvpadsbfjvajdfbs"
	tokenSign = "jbvirfbvabjdljkhnlnmcldksnclkdns"
	tokenTTL  = 12
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func GenerateUserToken(user structures.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Username,
	})
	return token.SignedString([]byte(tokenSign))
}

func ParseToken(authKey string) (int, error) {
	token, err := jwt.ParseWithClaims(authKey, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, errors.New("invalid token " + authKey)
		}

		return []byte(tokenSign), nil
	})

	if err != nil {
		logrus.Errorf("%v", err)
		return 0, errors.New("invalid token " + authKey)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("invalid token " + authKey)
	}
	logrus.Infof("%v", claims)
	return claims.UserId, nil
}

func GetUser(data requests.LoginRequest) (*structures.User, error) {
	user, err := db.GetUserRecord(data.Username)

	if err == sql.ErrNoRows {
		return nil, errors.New("username/password are not valid")
	}

	if err != nil {
		return nil, err
	}

	if validatePassword(data.Password, user.PasswordHash) {
		return user, nil
	}
	return nil, errors.New("username/password are not valid")
}

func validatePassword(password, hash string) bool {
	return generatePasswordHash(password) == hash
}

func CreateUser(data requests.RegisterRequest) (int, error) {
	user := structures.User{Username: data.Username}
	user.PasswordHash = generatePasswordHash(data.Password)
	return db.CreateUserRecord(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
