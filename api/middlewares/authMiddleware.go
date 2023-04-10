package middleware

import (
	"bcraft/api/errs"
	"bcraft/api/modules/auth"
	"bcraft/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			errs.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			return
		}

		uid, err := auth.ParseToken(headerParts[1])
		if err != nil {
			errs.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		user, err := db.GetUserRecordById(uid)

		if err != nil {
			errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.Set("user", user)
		c.Set("uid", uid)
	}
}
