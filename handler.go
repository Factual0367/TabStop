package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchRequest(c *gin.Context) {
	query := c.Query("query")
	Songs := GetSongList(query)
	AddGPDownloadLinks(Songs)
	SongsMap := make(map[int]Song)
	for index, item := range Songs {
		SongsMap[index] = item
	}
	c.IndentedJSON(http.StatusOK, SongsMap)
}
