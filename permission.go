package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Permission struct {
	Unit        string `json:"category" binding:"required"`
	Name        string `json:"cn" binding:"required"`
	Description string
	MType       []string
	UType       []string
}

func (p *Permission) rbacTypes() []string {
	return append(p.MType, p.UType...)
}

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
	var p Permission
	err := c.BindJSON(&p) // 会发送错误信息
	if err != nil {
		return
	}

	id, err := org.AddPermission(p.Name, p.Description, p.rbacTypes(), p.Unit == `department`)
	if err != nil {
		sendError(c, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

func permissionByID(c *gin.Context) {
	p, e := org.PermissionByID(c.Param(`id`))
	if e != nil {
		sendError(c, e)
		return
	}

	isUnit := p[`isUnit`].(bool)
	if isUnit {
		p[`category`] = `部门权限`
		p[`uType`] = p[`rbacType`]
	} else {
		p[`category`] = `员工权限`
		p[`mType`] = p[`rbacType`]
	}

	delete(p, `dn`)
	delete(p, `rbacType`)

	c.JSON(http.StatusOK, p)
}

func editPermission(c *gin.Context) {

	id := c.Param(`id`)
	var p Permission
	err := c.BindJSON(&p) // 会发送错误信息
	if err != nil {
		return
	}

	err = org.ModifyPermission(id, p.Name, p.Description, p.rbacTypes())
	if err != nil {
		sendError(c, err)
	} else {
		c.JSON(http.StatusOK, map[string]string{
			`id`: id,
		})
	}
}

// 删除逻辑
// 保证没人引用这个Permission
func delPermission(c *gin.Context) {
	id := c.Param(`id`)
	err := org.DelPermission(id)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		`id`: id,
	})
}
