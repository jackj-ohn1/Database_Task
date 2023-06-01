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
	bookHandler := handler.NewBookHandler(db)
	
	baseRoute.POST("/register", userHandler.Register)
	baseRoute.POST("/login", userHandler.Login)
	
	userRoute := baseRoute.Group("/user/book")
	
	// curl localhost:80/api/v1/library/user/book/whole?page=1&limit=10
	userRoute.GET("/whole", bookHandler.GetBooks)
	
	// curl -X POST -H "Content-Type:application/json" localhost:80/api/v1/library/user/book/borrow -d '{"user_id":"2021","book_id":1000000,"borrow_day":2}'
	userRoute.POST("/borrow", bookHandler.LendBook)
	
	// curl -X DELETE -H "Content-Type:application/json" localhost:80/api/v1/library/user/book/borrow -d '{"borrow_id":1000004}'
	userRoute.PUT("/borrow", userHandler.ReturnBook)
	
	// curl -X GET -H "Content-Type:application/json" localhost:80/api/v1/library/user/book/borrow/history -d '{"user_id":"2021"}'
	userRoute.GET("/borrow/history", userHandler.UserHistory)
	userRoute.GET("/borrow", userHandler.UserReservation)
	
	adminRoute := baseRoute.Group("/admin")
	// curl -X POST -H "Content-Type:application/json" localhost:80/api/v1/library/admin/book -d '{"user_id":"2023","book_name":"book1","book_author":"author1"}'
	adminRoute.POST("/book", bookHandler.AddBook)
	adminRoute.PUT("/book", bookHandler.DeleteBook)
	
	return engine
}
