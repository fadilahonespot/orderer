package controller

import (
	"net/http"

	"github.com/fadilahonespot/orderer/entity"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type defautUser struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &defautUser{db}
}

func (s *defautUser) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(400, response)
		return
	}

	pass, _ := hashPassword(req.Password)
	reqUser := entity.User{
		Email:    req.Email,
		Password: pass,
	}

	var userData entity.User
	s.db.Take(&userData, "email = ?", reqUser.Email)
	if userData.Email != "" {
		response := generateResponseError(http.StatusBadRequest, "Email already exists")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := s.db.Create(&reqUser).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	req.Password = pass
	response := generateResponseSuccess(http.StatusCreated, "Register Success", req)
	c.JSON(http.StatusCreated, response)
}

func (s *defautUser) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(400, response)
		return
	}

	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", req.Email)
	if reqUser.Email == "" {
		response := generateResponseError(http.StatusNotFound, "Email not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	err := comparePassword(reqUser.Password, req.Password)
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var dataResp LoginResponse
	token, err := createToken(cast.ToString(reqUser.ID))
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	dataResp.Token = token

	response := generateResponseSuccess(http.StatusOK, "Login Success", dataResp)
	c.JSON(http.StatusOK, response)
}
