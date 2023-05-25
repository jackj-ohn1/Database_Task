package route

import (
	"control/internal/handler"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func Router(db *sql.DB) *gin.Engine {
	engine := gin.Default()
	
	baseRoute := engine.Group("/api/v1/library")
	
	userHandler := handler.NewUserHandler(db)
	baseRoute.POST("/register", userHandler.Register)
	baseRoute.POST("/login", userHandler.Login)
	
	userRoute := engine.Group("/user/book")
	userRoute.GET("/")
	userRoute.POST("/borrow")
	userRoute.GET("/borrow")
	userRoute.GET("/borrow/history")
	
	adminRoute := engine.Group("/admin")
	adminRoute.POST("/book")
	adminRoute.DELETE("/book")
	
	return engine
}
