package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func fetchMembers(c *gin.Context) {

	_, pageSize, cookie := findPageControl(c)

	sr, err := org.Members(pageSize, cookie)
	if err != nil {
		sendError(c, err)
		return
	}

	sendResult(c, sr)
}

func mapToMemberInfo(body map[string]interface{}) map[string][]string {

	memberInfo := make(map[string][]string, 0)

	if name, ok := body[`name`].(string); ok {
		memberInfo[`name`] = []string{name}
	}
	if cn, ok := body[`cn`].(string); ok {
		memberInfo[`cn`] = []string{cn}
	}
	if email, ok := body[`email`].(string); ok {
		memberInfo[`email`] = []string{email}
	}
	if telephoneNumber, ok := body[`telephoneNumber`].(string); ok {
		memberInfo[`telephoneNumber`] = []string{telephoneNumber}
	}
	if title, ok := body[`title`].(string); ok {
		memberInfo[`title`] = []string{title}
	}
	if unitID, ok := body[`unitID`].(string); ok {
		memberInfo[`unitID`] = []string{unitID}
	}
	if rbacRoles, ok := body[`rbacRole`].([]interface{}); ok {
		var roles []string
		for _, r := range rbacRoles {
			roles = append(roles, r.(string))
		}
		memberInfo[`rbacRole`] = roles
	}
	if rbacTypes, ok := body[`rbacType`].([]interface{}); ok {
		var types []string
		for _, t := range rbacTypes {
			types = append(types, t.(string))
		}
		memberInfo[`rbacType`] = types
	}

	return memberInfo
}

func createMember(c *gin.Context) {

	var body map[string]interface{}
	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}

	memberInfo := mapToMemberInfo(body)

	m := md5.New()
	m.Write([]byte(`123456`))
	pwd := m.Sum(nil)

	memberInfo[`userPassword`] = []string{fmt.Sprintf(`{MD5}%s`, hex.EncodeToString(pwd))}

	id, err := org.AddMember(memberInfo)
	if err != nil {
		sendError(c, err)
		return
	}

	thirdPwd := newPassword()
	r, err := em.RegisterSignelUser(id, thirdPwd)
	if err != nil {
		log.WithField(`resource`, `member`).Warn(fmt.Sprintf(`register user failed %v`, err))
	}
	if !r {
		log.WithField(`resource`, `member`).Warn(`register user failed with no reason`)
	}
	// 将用户名和密码更新到ldapserver
	err = org.ModifyMember(id, map[string][]string{
		`thirdAccount`:  []string{id},
		`thirdPassword`: []string{thirdPwd},
	})
	if err != nil {
		log.WithField(`resource`, `memeber`).Warn(fmt.Sprintf(`modify member third account info occured error %v`, err))
	}

	c.JSON(http.StatusOK, map[string]string{`id`: id})
	// 去环信注册
}
