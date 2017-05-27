package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchPermissionByRoleUPID(c *gin.Context) {
	fetchPermissionByRole(true, c)
}

func fetchPermissionByRolePPID(c *gin.Context) {
	fetchPermissionByRole(false, c)
}

func fetchPermissionByRole(isUnit bool, c *gin.Context) {

	rid := c.Query(`role_id`) // 有可能是查找某个角色包含的所有类型
	r, e := org.RoleByID(rid)
	if e != nil {
		sendError(c, e)
		return
	}

	ids := r[`ppid`].([]string)
	if isUnit {
		ids = r[`upid`].([]string)
	}

	sr, e := org.PermissionByIDs(ids)
	if e != nil {
		sendError(c, e)
		return
	}
	sendResult(c, sr)
}

func fetchPermissions(c *gin.Context) {

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

func createPermission(c *gin.Context) {
	var body map[string]interface{}
	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}

	category := body[`category`]

	isUnit := category == `1`
	name := body[`name`].(string)
	desc := body[`description`].(string)
	types := body[`rbacType`].([]string)

	id, err := org.AddPermission(name, desc, types, isUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

func permissionByID(c *gin.Context) {
	p, e := org.PermissionByID(c.Param(`id`))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			`err`: e,
		})
		return
	}
	// 做点小处理 对前端友好
	if p[`isUnit`].(bool) {
		p[`category`] = department
	} else {
		p[`category`] = member
	}
	delete(p, `isUnit`)
	c.JSON(http.StatusOK, p)
}

func editPermission(c *gin.Context) {

	id := c.Param(`id`)
	var body map[string]interface{}
	e := c.BindJSON(&body)
	if e != nil {
		return
	}

	name := body[`name`].(string)
	desc := body[`description`].(string)
	types := body[`rbacType`].([]string)
	isUnit := body[`category`].(string) == department

	e = org.ModifyPermission(id, name, desc, types, isUnit)
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

func delPermission(c *gin.Context) {

	p, e := org.PermissionByID(c.Param(`id`))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
		return
	}

	e = org.DelPermission(p[`id`].(string), p[`isUnit`].(bool))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, p)
}
