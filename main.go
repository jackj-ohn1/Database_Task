package main

import (
	"control/internal/repository"
	"control/route"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("debug")
	db := repository.Init("SA", "")
	engine := route.Router(db)
	engine.Run(":80")
}
