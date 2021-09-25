package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrBadRequest(g *gin.Context, err error) {
	g.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": false,
		"code":    http.StatusBadRequest,
		"error":   err.Error(),
	})
}

func ErrInternalServerError(g *gin.Context, err error) {
	g.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"error":   err.Error(),
	})
}

func ErrUnauthorized(g *gin.Context, err error) {
	g.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"error":   err.Error(),
	})
}

func Success(g *gin.Context, data interface{}) {
	g.AbortWithStatusJSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data,
	})
}
