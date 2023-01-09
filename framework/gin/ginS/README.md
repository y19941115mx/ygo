# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/y19941115mx/ygo/framework/gin"
	"github.com/y19941115mx/ygo/framework/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
