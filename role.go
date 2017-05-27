package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchRoles(c *gin.Context) {

	// pid := c.Query(`permission_id`) // 有可能是查找某个权限包含的所有类型
	//
	// if len(pid) > 0 {
	//
	// 	ts, e := org.TypeByPermissionID(pid)
	// 	if e != nil {
	// 		c.JSON(http.StatusInternalServerError, map[string]interface{}{
	// 			`err`: e.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, map[string]interface{}{
	// 		`total`: 1,
	// 		`data`:  ts,
	// 		`page`:  1,
	// 	})
	// 	return
	// }

	_, pageSize, cookie := findPageControl(c)

	r, e := org.Roles(pageSize, cookie)
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
	}
	sendResult(c, r)
}
