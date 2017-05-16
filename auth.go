package main

import (
	"time"

	"github.com/gin-gonic/gin"
	jwt "gopkg.in/appleboy/gin-jwt.v2"
)

// ClientJWTMiddleware ...
func ClientJWTMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      `Dolores Client Zone`,
		Key:        []byte(`secret key`),
		Timeout:    24 * time.Hour,     // token 有效期一天
		MaxRefresh: 7 * 24 * time.Hour, // 一周以内可以刷新
		Authenticator: func(userID string, password string, c *gin.Context) (string, bool) {
			id, err := org.AuthMember(userID, password)
			return id, err == nil
		},
		TokenHeadName: `Dolores`,
	}

	return authMiddleware
}
