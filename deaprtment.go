package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchDepartment(c *gin.Context) {

	r, e := org.AllUnit()
	if e != nil {
		sendError(c, e)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		`data`:  r,
		`total`: len(r),
	})
}

func createDepartment(c *gin.Context) {

	var body map[string]string

	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}

	parentID := body[`parentID`]

	info := map[string][]string{
		`ou`:          []string{body[`cn`]},
		`description`: []string{body[`description`]},
		`rbacType`:    []string{body[`rbacType`]},
	}

	id, err := org.AddUnit(parentID, info)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{`id`: id})
}

func departmentByID(c *gin.Context) {

	unit, err := org.UnitByID(c.Param(`id`))
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		`data`: unit,
	})
}
