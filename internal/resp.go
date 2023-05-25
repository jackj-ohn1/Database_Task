package internal

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int    `json:"code"`
}

func (r *Response) SomeError(c *gin.Context, err error) {
	log.Println("some err happened:", err)
	c.JSON(http.StatusOK, r)
}

func (r *Response) InternalError(c *gin.Context, err error) {
	log.Println("internal server error:", err)
	c.JSON(http.StatusInternalServerError, r)
	c.Abort()
}

func (r *Response) Success(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}
