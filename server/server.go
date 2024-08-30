package server

import (
	"github.com/gin-gonic/gin"
	"github.com/onurhanak/songsterrapi"
)

func Start() {
	router := gin.Default()
	router.GET("/search", songsterrapi.SearchRequest)
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
