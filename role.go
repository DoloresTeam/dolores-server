package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Role struct {
	Name        string `json:"cn"`
	Description string
	UPS         []string `json:"upid"`
	PPS         []string `json:"ppid"`
}

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

func createRole(c *gin.Context) {
	var r Role
	err := c.BindJSON(&r) // 会发送错误信息
	if err != nil {
		return
	}

	id, err := org.AddRole(r.Name, r.Description, r.UPS, r.PPS)
	if err != nil {
		sendError(c, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

func roleByID(c *gin.Context) {
	r, e := org.RoleByID(c.Param(`id`))
	if e != nil {
		sendError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}

func editRole(c *gin.Context) {

	id := c.Param(`id`)
	var r Role
	err := c.BindJSON(&r) // 会发送错误信息
	if err != nil {
		return
	}

	err = org.ModifyRole(id, r.Name, r.Description, r.UPS, r.PPS)
	if err != nil {
		sendError(c, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

// 删除逻辑
func delRole(c *gin.Context) {
	id := c.Param(`id`)
	e := org.DelRole(id)
	if e != nil {
		sendError(c, e)
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		`id`: id,
	})
}
