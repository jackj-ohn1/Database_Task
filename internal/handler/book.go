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
	"strconv"
)

type bookHandler struct {
	database repository.Database
}

func NewBookHandler(db *sql.DB) *bookHandler {
	return &bookHandler{
		repository.NewDatabase(db),
	}
}

func (b *bookHandler) GetBooks(c *gin.Context) {
	pageStr, limitStr := c.DefaultQuery("page", "1"), c.DefaultQuery("limit", "10")
	
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	data, err := b.database.GetBooks(page, limit)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	resp := &internal.Response{
		Message: "获取成功",
		Code:    http.StatusOK,
		Data:    data,
	}
	resp.Success(c)
}

type lendBookRequest struct {
	UserId    string `json:"user_id"`
	BookId    int    `json:"book_id"`
	BorrowDay int    `json:"borrow_day"`
}

func (b *bookHandler) LendBook(c *gin.Context) {
	var req lendBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	err := b.database.LendBook(&model.Borrow{
		UserId: req.UserId,
		BookId: req.BookId,
	}, req.BorrowDay)
	if err != nil {
		if errors.Is(errno.ErrUsedBook, err) {
			resp := &internal.Response{
				Message: err.Error(),
				Code:    http.StatusOK,
			}
			resp.Success(c)
			return
		}
		
		resp := &internal.Response{
			Message: errno.ErrInternalServerErr.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	resp := &internal.Response{
		Message: "操作成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}
