package main

import (
	"net/http"

	"qiniupkg.com/api.v7/kodo"

	"github.com/gin-gonic/gin"
)

func profile(c *gin.Context) {

	id, _ := c.Get(`userID`)

	member, err := org.MemberByID(id.(string), false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			`err`: err.Error(),
		})
	} else {
		delete(member, `dn`)

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

func qiniuUploadToken(c *gin.Context) {

	k := kodo.New(0, nil)
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: `dolores`,
		//设置Token过期时间  5 分钟
		Expires: 3600 * 5,
	}
	// 生成一个上传token
	token := k.MakeUptoken(policy)

	if len(token) != 0 {
		c.JSON(http.StatusOK, map[string]string{
			`token`: token,
		})
	} else {
		c.JSON(http.StatusInternalServerError, map[string]string{
			`err`: `can't construct token`,
		})
	}

}
