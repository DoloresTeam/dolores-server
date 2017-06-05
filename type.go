package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func fetchTypes(c *gin.Context) {
	pid := c.Query(`permission_id`) // 有可能是查找某个权限包含的所有类型
	if len(pid) > 0 {
		ts, e := org.TypeByPermissionID(pid)
		if e != nil {
			sendError(c, e)
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			`total`: len(ts),
			`data`:  ts,
			`page`:  1,
		})
		return
	}
	idStr := c.Query(`ids`) // 通过一组ID 获取type
	if len(idStr) > 0 {
		ts, e := org.TypeByIDs(strings.Split(idStr, `,`))
		if e != nil {
			sendError(c, e)
			return
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			`total`: len(ts),
			`data`:  ts,
			`page`:  1,
		})
		return
	}

	isUnit := c.Query(`isUnit`) == `true` // 需要判断是查找部门类型，还是员工类型

	_, pageSize, cookie := findPageControl(c)

	r, e := org.Types(isUnit, pageSize, cookie)
	if e != nil {
		sendError(c, e)
		return
	}
	sendResult(c, r)
}

func createType(c *gin.Context) {

	var body map[string]string

	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}

	category := body[`category`]
	isUnit := category == department

	id, err := org.AddType(body[`cn`], body[`description`], isUnit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{`id`: id}) // 将id 返回
}

func typeByID(c *gin.Context) {

	t, e := org.TypeByID(c.Param(`id`))
	if e != nil {
		sendError(c, e)
		return
	}
	// 做点小处理 对前端友好
	if t[`isUnit`].(bool) {
		t[`category`] = department
	} else {
		t[`category`] = member
	}
	// delete(t, `isUnit`)
	c.JSON(http.StatusOK, t)
}

func editType(c *gin.Context) {

	id := c.Param(`id`)
	var body map[string]interface{}
	e := c.BindJSON(&body)
	if e != nil {
		return
	}
	e = org.ModifyType(id, body[`cn`].(string), body[`description`].(string), body[`category`].(string) == department)
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

// Type 删除逻辑
// 如果有权限包含该Type 则提示，不能删除成功
// 如果有员工属于该Type则提示，不能删除成功
func delType(c *gin.Context) {

	id := c.Param(`id`)
	t, e := org.TypeByID(id)
	if e != nil {
		sendError(c, e)
		return
	}

	e = org.DelType(t[`id`].(string), t[`isUnit`].(bool))
	if e != nil {
		sendError(c, e)
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{`data`: t})
	}
}
