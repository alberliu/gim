# Favicon Gin's middleware

Gin middleware to support favicon.

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/thinkerou/favicon
```

Import it in your code:

```go
import "github.com/thinkerou/favicon"
```

### Canonical example:

```go
package main
            
import (
    "github.com/gin-gonic/gin"
    "github.com/thinkerou/favicon"
)
            
func main() {
    r := gin.Default()
    r.Use(favicon.New("./favicon.ico")) // set favicon middleware 

    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "Hello favicon.")
    })

    r.Run(":8080")
}
```
