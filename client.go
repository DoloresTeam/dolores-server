package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func profile(c *gin.Context) {

	id, _ := c.Get(`userID`)

	member, err := org.MemberByID(id.(string), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			`err`: err.Error(),
		})
	} else {
		delete(member, `rbacRole`)
		delete(member, `rbacType`)

		c.JSON(http.StatusOK, member)
	}
}

func updateAvatarURL(c *gin.Context) {

	id, _ := c.Get(`userID`)

	body := make(map[string]interface{}, 0)

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			`err`: `缺少avatarURL 参数`,
		})
	} else {
		err = org.ModifyMember(id.(string), map[string][]string{
			`labeledURI`: []string{body[`avatarURL`].(string)},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				`err`: err.Error(),
			})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func organizationMap(c *gin.Context) {

	id, _ := c.Get(`userID`)

	members, err := org.OrganizationMemberByMemberID(id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: err.Error(),
		})
		return
	}
	departments, err := org.OrganizationUnitByMemberID(id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		`members`:     members,
		`departments`: departments,
		`version`:     1,
	})
}
