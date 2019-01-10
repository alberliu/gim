// Package autotls support Let's Encrypt for a Go server application.
//
// 	package main
//
// 	import (
// 		"log"
//
// 		"github.com/gin-gonic/autotls"
// 		"github.com/gin-gonic/gin"
// 	)
//
// 	func main() {
// 		r := gin.Default()
//
// 		// Ping handler
// 		r.GET("/ping", func(c *gin.Context) {
// 			c.String(200, "pong")
// 		})
//
// 		log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
// 	}
//
package autotls
