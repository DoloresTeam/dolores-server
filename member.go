package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/fogleman/gg"
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
	if gender, ok := body[`gender`].(string); ok {
		memberInfo[`gender`] = []string{gender}
	}
	if priority, ok := body[`priority`].(string); ok {
		memberInfo[`priority`] = []string{priority}
	}
	if telephoneNumber, ok := body[`telephoneNumber`].(string); ok {
		memberInfo[`telephoneNumber`] = []string{telephoneNumber}
	}
	if title, ok := body[`title`].(string); ok {
		memberInfo[`title`] = []string{title}
	}
	if unitIDs, ok := body[`unitID`].([]interface{}); ok {
		var ids []string
		for _, id := range unitIDs {
			ids = append(ids, id.(string))
		}
		memberInfo[`unitID`] = ids
	}
	if rbacType, ok := body[`rbacType`].(string); ok {
		memberInfo[`rbacType`] = []string{rbacType}
	}
	if rbacRoles, ok := body[`rbacRole`].([]interface{}); ok {
		var roles []string
		for _, r := range rbacRoles {
			roles = append(roles, r.(string))
		}
		memberInfo[`rbacRole`] = roles
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

	// 去环信注册
	err = em.RegisterSignelUser(id, thirdPwd)
	if err != nil {
		log.WithField(`resource`, `member`).Warn(fmt.Sprintf(`register user failed %v`, err))
	}
	// 将用户名和密码更新到ldapserver
	err = org.ModifyMember(id, map[string][]string{
		`thirdAccount`:  []string{id},
		`thirdPassword`: []string{thirdPwd},
	})
	if err != nil {
		log.WithField(`resource`, `memeber`).Warn(fmt.Sprintf(`modify member third account info occured error %v`, err))
	}

	// go generatorAvatar(id, body[`name`].(string))
	c.JSON(http.StatusOK, map[string]string{`id`: id})
}

func generatorAvatar(id, name string) {

	dc := gg.NewContext(210, 210)
	dc.DrawCircle(105, 105, 105)
	dc.SetHexColor(`#03A9F4`)
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace(`../STHeiti Medium.ttc`, 80)
	dc.DrawStringAnchored(name, 105, 105, 0.5, 0.5)
	dc.Stroke()

	dc.SavePNG("out.png")

}

func editMember(c *gin.Context) {
	var body map[string]interface{}
	err := c.BindJSON(&body) // 会发送错误信息
	if err != nil {
		return
	}

	memberInfo := mapToMemberInfo(body)

	id := c.Param(`id`)
	err = org.ModifyMember(id, memberInfo)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{`id`: id})
}

func memberByID(c *gin.Context) {
	ms, err := org.MemberByID(c.Param(`id`), true, false)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, ms)
}

func delMember(c *gin.Context) {

	id := c.Param(`id`)

	err := org.DelMember(id)
	if err != nil {
		sendError(c, err)
		return
	}
	// 去环信删除
	err = em.DeleteUser(id)
	if err != nil {
		log.WithField(`resource`, `member`).Warn(fmt.Sprintf(`delete user[id:%v] failed %v`, id, err))
	}

	c.JSON(http.StatusOK, map[string]string{`id`: id})
}
