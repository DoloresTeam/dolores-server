package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func fetchDepartment(c *gin.Context) {
	idStr := c.Query(`ids`)
	var us []map[string]interface{}
	var err error
	if len(idStr) > 0 {
		us, err = org.UnitByIDs(strings.Split(idStr, `,`))
	} else {
		us, err = org.AllUnit()
	}

	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		`data`:  us,
		`total`: len(us),
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

func editDepartment(c *gin.Context) {

	var body map[string]string
	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}
	info := map[string][]string{
		`ou`:          []string{body[`cn`]},
		`description`: []string{body[`description`]},
		`rbacType`:    []string{body[`rbacType`]},
	}

	id := c.Param(`id`)
	err = org.ModifyUnit(id, info)
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

func delDepartment(c *gin.Context) {
	id := c.Param(`id`)
	err := org.DelUnit(id)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		`id`: id,
	})
}
