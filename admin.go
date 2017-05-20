package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DoloresTeam/organization"
	"github.com/gin-gonic/gin"
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

func types(c *gin.Context) {

	pid := c.Query(`permission_id`)
	if len(pid) > 0 {

		ts, e := org.TypeByPermissionID(pid)
		if e != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				`err`: e.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			`total`: 1,
			`data`:  ts,
			`page`:  1,
		})
		return
	}

	isUnit := c.Query(`isUnit`) == `true`

	_, pageSize, cookie := findPageControl(c)

	r, e := org.Types(isUnit, pageSize, cookie)
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
	}
	sendResult(c, r)
}

func createType(c *gin.Context) {

	var body map[string]string
	err := c.BindJSON(&body)
	if err != nil {
		return
	}
	category := body[`category`]

	isUnit := category == `1`
	id, err := org.AddType(body[`cn`], body[`description`], isUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{`id`: id})
}

func permissions(c *gin.Context) {

	isUnit := c.Query(`isUnit`) == `true`

	_, pageSize, cookie := findPageControl(c)

	r, e := org.Permissions(isUnit, pageSize, cookie)
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
	}
	sendResult(c, r)
}

func findPageControl(c *gin.Context) (uint32, uint32, []byte) {

	p, e := strconv.Atoi(c.DefaultQuery(`pageSize`, `25`))
	if e != nil {
		p = 25
	}
	pageSize := uint32(p)

	p0, e := strconv.Atoi(c.DefaultQuery(`page`, `1`))
	if e != nil {
		p0 = 1
	}
	page := uint32(p0)
	if page == 0 {
		page = 1
	}

	var cookie []byte
	if len(c.Query(`pageCookie`)) > 0 {
		cookie = []byte(cookie)
	}

	return page, pageSize, cookie
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
