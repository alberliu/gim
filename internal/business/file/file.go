package file

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RunFileServer() {
	router := gin.Default()
	router.Static("/file", "/data/file")

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		name := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := "/data/file/" + name
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		baseURL := fmt.Sprintf("http://%s/file/", c.Request.Host)
		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			Data:    map[string]string{"url": baseURL + name},
		})
	})
	err := router.Run(":8081")
	if err != nil {
		slog.Error("Run error", "error", err)
		panic(err)
	}
}
