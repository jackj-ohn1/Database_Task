package handler

import (
	"control/internal"
	"control/internal/errno"
	"control/internal/repository"
	"control/internal/repository/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type userHandler struct {
	database repository.Database
}

func NewUserHandler(db *sql.DB) *userHandler {
	return &userHandler{
		repository.NewDatabase(db),
	}
}

type userRegisterRequest struct {
	UserId       string `json:"user_id" binding:"required,min=8,max=10"`
	UserPassword string `json:"user_password" binding:"required,min=8,max=16"`
	UserName     string `json:"user_name" binding:"required"`
	Invitation   string `json:"invitation"`
}

func (u *userHandler) Register(c *gin.Context) {
	var req userRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	err := u.database.UserRegister(&model.User{
		UserId:       req.UserId,
		UserPassword: req.UserPassword,
		UserName:     req.UserName,
	})
	if err != nil {
		if errors.Is(err, errno.ErrUserExist) {
			resp := &internal.Response{
				Message: err.Error(),
				Code:    http.StatusOK,
			}
			resp.Success(c)
			return
		}
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	resp := &internal.Response{
		Message: "注册成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}

type userLoginRequest struct {
	UserId       string `json:"user_id" binding:"required,min=8,max=10"`
	UserPassword string `json:"user_password" binding:"required,min=8,max=16"`
}

func (u *userHandler) Login(c *gin.Context) {
	var req userLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	err := u.database.UserLogin(&model.User{
		UserId:       req.UserId,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		if errors.Is(err, errno.ErrWrongPassword) {
			resp := &internal.Response{
				Message: err.Error(),
				Code:    http.StatusOK,
			}
			resp.Success(c)
			return
		}
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	resp := &internal.Response{
		Message: "登录成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}
