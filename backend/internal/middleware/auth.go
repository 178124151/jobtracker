package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1001,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// 检查Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    1001,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		token := parts[1]
		
		// TODO: 验证JWT token
		// 临时实现：解析token获取user_id
		_ = token // 将在JWT实现中使用
		userID := "temp-user-id" // 实际应该从JWT解析
		
		c.Set("user_id", userID)
		c.Next()
	}
}