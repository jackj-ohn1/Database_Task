package handler

import (
	"control/internal"
	"control/internal/errno"
	"control/internal/repository"
	"control/internal/repository/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type bookHandler struct {
	database repository.Database
	cache    *lru.Cache[int, *model.Book]
	size     int
	isNew    bool
}

func NewBookHandler(db *sql.DB) *bookHandler {
	b := &bookHandler{
		size:     50,
		database: repository.NewDatabase(db),
		isNew:    true,
	}
	cache, _ := lru.New[int, *model.Book](b.size)
	b.cache = cache
	return b
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
	
	if !b.isNew {
		data := make([]*model.Book, 0, limit)
		keys := b.cache.Keys()
		for i := (page - 1) * limit; i < page*limit && i < len(keys); i++ {
			value, ok := b.cache.Get(keys[i])
			if !ok {
				continue
			}
			data = append(data, value)
		}
		
		resp := &internal.Response{
			Message: "获取成功",
			Code:    http.StatusOK,
			Data:    data,
		}
		resp.Success(c)
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
	
	b.isNew = false
	if len(data) > b.size {
		b.cache.Resize(len(data) * 2)
		b.size = len(data) * 2
	}
	
	var bookSlice model.BookSlice = data
	
	for _, v := range bookSlice {
		b.cache.ContainsOrAdd(v.BookId, v)
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
		} else if errors.Is(errno.ErrNoLeftResource, err) {
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
	
	b.cache.Get(req.BookId)
	
	resp := &internal.Response{
		Message: "操作成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}

type addBookRequest struct {
	UserId          string    `json:"user_id" binding:"required"`
	BookName        string    `json:"book_name" binding:"required"`
	BookAuthor      string    `json:"book_author"`
	BookPublishTime time.Time `json:"book_publish_time"`
}

func (b *bookHandler) AddBook(c *gin.Context) {
	var req addBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	isAdmin, err := b.database.IsAdmin(req.UserId)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	} else if !isAdmin {
		resp := &internal.Response{
			Message: errno.ErrNoPower.Error(),
			Code:    http.StatusUnauthorized,
		}
		resp.Success(c)
		return
	}
	
	err = b.database.AddBook(&model.Book{
		BookName:   req.BookName,
		BookAuthor: req.BookAuthor,
		BookPublishedTime: sql.NullTime{
			Time:  req.BookPublishTime,
			Valid: true,
		},
	})
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	id, err := b.database.GetBookId(req.BookName)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	b.cache.ContainsOrAdd(id, &model.Book{
		BookId:     id,
		BookName:   req.BookName,
		BookAuthor: req.BookAuthor,
		BookPublishedTime: sql.NullTime{
			Time:  req.BookPublishTime,
			Valid: true,
		},
	})
	if b.cache.Len() >= b.size {
		b.cache.Resize(b.size * 2)
		b.size = b.size * 2
	}
	
	resp := &internal.Response{
		Message: "添加成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}

type deleteBookRequest struct {
	UserId string `json:"user_id" binding:"required"`
	BookId int    `json:"book_id" binding:"required"`
}

func (b *bookHandler) DeleteBook(c *gin.Context) {
	var req deleteBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		resp.SomeError(c, err)
		return
	}
	
	isAdmin, err := b.database.IsAdmin(req.UserId)
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	} else if !isAdmin {
		resp := &internal.Response{
			Message: errno.ErrNoPower.Error(),
			Code:    http.StatusUnauthorized,
		}
		resp.Success(c)
		return
	}
	
	err = b.database.DeleteBook(&model.Book{
		BookId: req.BookId,
	})
	if err != nil {
		resp := &internal.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		}
		resp.InternalError(c, err)
		return
	}
	
	b.cache.Remove(req.BookId)
	
	resp := &internal.Response{
		Message: "删除成功",
		Code:    http.StatusOK,
	}
	resp.Success(c)
}
