package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DoloresTeam/organization"
	"github.com/gin-gonic/gin"
)

const (
	department = `Department`
	member     = `Member`
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func findPageControl(c *gin.Context) (page uint32, pageSize uint32, cookie []byte) {

	p, e := strconv.Atoi(c.DefaultQuery(`pageSize`, `25`))
	if e != nil {
		p = 25
	}
	pageSize = uint32(p)

	p0, e := strconv.Atoi(c.DefaultQuery(`page`, `1`))
	if e != nil {
		p0 = 1
	}
	page = uint32(p0)
	if page == 0 {
		page = 1
	}

	if len(c.Query(`pageCookie`)) > 0 {
		cookie = []byte(cookie)
	}

	return
}

func sendResult(c *gin.Context, r *organization.SearchResult) {

	page, pageSize, _ := findPageControl(c)

	total := (page-1)*pageSize + uint32(len(r.Data))
	if r.Cookie != nil {
		total += pageSize
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		`pageSize`: r.Size,
		`total`:    total,
		`cookie`:   string(r.Cookie),
		`data`:     r.Data,
		`page`:     page,
	})
}

func sendError(c *gin.Context, e error) {
	c.JSON(http.StatusInternalServerError, map[string]string{
		`message`: e.Error(),
	})
}
