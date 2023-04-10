package auth

import (
	"bcraft/api/errs"
	request "bcraft/api/structures/requests"
	response "bcraft/api/structures/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct{}

// @Summary Login
// @Tags auth
// @Description Get JWTtoken for authorization header
// @ID login
// @Accept json
// @Produce json
// @Param Credentials body request.LoginRequest true "Credentials"
// @Success 200 {object} response.LoginResponce
// @Failure 400
// @Failure 500
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var body request.LoginRequest

	if err := c.BindJSON(&body); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := GetUser(body)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		errs.NewErrorResponse(c, http.StatusUnauthorized, "wrong credentials")
		return
	}
	token, err := GenerateUserToken(*user)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.LoginResponce{Token: token})
}

// @Summary Register a user
// @Tags auth
// @Description Create user
// @ID register
// @Accept json
// @Produce json
// @Param Credentials body request.RegisterRequest true "Credentials"
// @Success 200 {object} response.SuccessResponse
// @Failure 400
// @Failure 500
// @Router /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var body request.RegisterRequest

	if err := c.BindJSON(&body); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err := CreateUser(body)

	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{Success: true})
}
