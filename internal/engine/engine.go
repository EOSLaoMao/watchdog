package engine

import "github.com/gin-gonic/gin"

var e = gin.Default()

func E() *gin.Engine {
	return e
}
