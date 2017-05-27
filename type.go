package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchTypes(c *gin.Context) {

	pid := c.Query(`permission_id`) // 有可能是查找某个权限包含的所有类型

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

	isUnit := c.Query(`isUnit`) == `true` // 需要判断是查找部门类型，还是员工类型

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

	err := c.BindJSON(&body) // 会发送错误信息
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

	c.JSON(http.StatusOK, map[string]string{`id`: id}) // 将id 返回
}

func typeByID(c *gin.Context) {

	t, e := org.TypeByID(c.Param(`id`))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			`err`: e,
		})
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

func delType(c *gin.Context) {

	t, e := org.TypeByID(c.Param(`id`))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
		return
	}
	e = org.DelType(t[`id`].(string), t[`isUnit`].(bool))
	if e != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: e.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{`data`: t})
	}
}
