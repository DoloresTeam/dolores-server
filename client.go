package main

import (
	"errors"
	"net/http"

	"qiniupkg.com/api.v7/kodo"

	"github.com/gin-gonic/gin"
)

// ModifyPassword ...
type ModifyPassword struct {
	OriginalPassword string `json:"originalPassword"`
	NewPassword      string `json:"newPassword"`
}

func profile(c *gin.Context) {

	id, _ := c.Get(`userID`)

	member, err := org.MemberByID(id.(string), false, true)
	if err != nil {
		sendError(c, err)
	} else {
		delete(member, `dn`)

		c.JSON(http.StatusOK, member)
	}
}

func basicProfiles(c *gin.Context) {
	ids := c.QueryArray(`id[]`)
	members, err := org.MemberByIDs(ids, false, false)
	if err != nil {
		sendError(c, err)
		return
	}
	var result []map[string]string
	for _, m := range members {
		result = append(result, map[string]string{
			`id`:         m[`id`].(string),
			`name`:       m[`name`].(string),
			`labeledURI`: m[`labeledURI`].(string),
		})
	}
	c.JSON(http.StatusOK, result)
}

func updateAvatarURL(c *gin.Context) {

	id, _ := c.Get(`userID`)

	body := make(map[string]interface{}, 0)

	err := c.BindJSON(&body)
	if err != nil {
		sendError(c, errors.New(`缺少avatarURL 参数`))
	} else {
		err = org.ModifyMember(id.(string), map[string][]string{
			`labeledURI`: []string{body[`avatarURL`].(string)},
		})
		if err != nil {
			sendError(c, err)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func modifyPassword(c *gin.Context) {
	id, _ := c.Get(`userID`)

	var mp ModifyPassword
	err := c.BindJSON(&mp)
	if err != nil {
		return
	}

	err = org.ModifyPassword(id.(string), mp.OriginalPassword, mp.NewPassword)
	if err != nil {
		sendError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func organizationMap(c *gin.Context) {
	id, _ := c.Get(`userID`)
	departments, members, version, err := org.OrganizationView(id.(string))
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		`departments`: departments,
		`members`:     members,
		`version`:     version,
	})
}

func syncOrganization(c *gin.Context) {
	id, _ := c.Get(`userID`)
	version := c.Param(`version`)
	if !org.IsValidVersion(version) {
		c.JSON(http.StatusOK, map[string]interface{}{
			`needRefetchOrganization`: true,
			`version`:                 version,
		})
		return
	}
	version, json, err := org.GenerateChangeLogFromVersion(version, id.(string))
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		`version`: version,
		`logs`:    json,
		`needRefetchOrganization`: false,
	})
}

func qiniuUploadToken(c *gin.Context) {

	k := kodo.New(0, nil)
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: `dolores`,
		//设置Token过期时间  5 分钟
		Expires: 3600 * 5,
	}
	// 生成一个上传token
	token, err := k.MakeUptokenWithSafe(policy)
	if err != nil {
		sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		`token`: token,
	})
}
