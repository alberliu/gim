package main

import (
	"gim/pkg/logger"
	"gim/pkg/util"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const baseUrl = "http://111.229.238.28:8085/file/"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// 初始化日志
	logger.Init()
	router := gin.Default()
	router.Static("/file", "/root/file")

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusOK, Response{Code: 1001, Message: err.Error()})
			return
		}

		filenames := strings.Split(file.Filename, ".")
		name := strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + util.RandString(30) + "." + filenames[len(filenames)-1]
		filePath := "/root/file/" + name
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			c.JSON(http.StatusOK, Response{Code: 1001, Message: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "success",
			Data:    map[string]string{"url": baseUrl + name},
		})
	})
	router.Run(":8085")
}
