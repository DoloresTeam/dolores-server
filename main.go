package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"qiniupkg.com/api.v7/conf"

	yaml "gopkg.in/yaml.v2"

	"github.com/DoloresTeam/organization"
	"github.com/gin-gonic/gin"
)

// Config ...
type Config struct {
	Host    string
	Port    int
	RootDN  string
	RootPWD string
	Subffix string
}

var config = Config{}
var org *organization.Organization

func main() {

	// 配置七牛的 ak 和 sk
	conf.ACCESS_KEY = `e1IUbuH8t2D-L9M3s9UmddOT_3YU7xk0VnB1Ws-8`
	conf.SECRET_KEY = `PTnbxt20SV47kkA-viiyfBrIdVlM4sqLDpv_wYJN`

	configFilePath := flag.String(`path`, `./conf.yaml`, `配置文件路径`)

	flag.Parse()

	_, err := os.Stat(*configFilePath)
	if !(err == nil || os.IsExist(err)) {
		panic(`配置文件不存在`)
	}
	data, _ := ioutil.ReadFile(*configFilePath)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(`配置文件不正确`)
	}

	_org, err := organization.NewOrganizationWithSimpleBind(config.Subffix,
		config.Host,
		config.RootDN,
		config.RootPWD,
		config.Port)
	if err != nil {
		panic(err.Error())
	}
	org = _org

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware()) // 允许跨站请求

	clientAuth := ClientJWTMiddleware()

	r.POST(`/login`, clientAuth.LoginHandler)

	auth := r.Group(`/api/v1`, clientAuth.MiddlewareFunc())
	{
		auth.GET(`/refresh_token`, clientAuth.RefreshHandler)
		auth.GET(`/upload_token`, qiniuUploadToken)
		auth.GET(`/profile`, profile)
		auth.POST(`/update_avatar`, updateAvatarURL)
		auth.GET(`/organization`, organizationMap)
	}

	// adminAuth := ServerJWTMiddleware()

	// r.POST(`/admin/login`, adminAuth.LoginHandler)

	admin := r.Group(`/admin/v1`)

	{
		admin.GET(`type`, types)
		admin.GET(`type/:id`, func(c *gin.Context) {
			t, e := org.TypeByID(c.Param(`id`))
			if e != nil {
				c.JSON(http.StatusInternalServerError, map[string]interface{}{
					`err`: e,
				})
				return
			}
			if t[`isUnit`].(bool) {
				t[`category`] = `Department`
			} else {
				t[`category`] = `Member`
			}
			delete(t, `isUnit`)
			c.JSON(http.StatusOK, t)
		})
		admin.PUT(`type/:id`, func(c *gin.Context) {
			id := c.Param(`id`)
			var body map[string]string
			e := c.BindJSON(&body)
			if e != nil {
				return
			}
			e = org.ModifyType(id, body[`cn`], body[`description`], body[`category`] == `Department`)
			if e != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{
					`err`: e.Error(),
				})
			} else {
				c.JSON(http.StatusOK, map[string]string{
					`id`: id,
				})
			}
		})
		admin.DELETE(`type/:id`, func(c *gin.Context) {
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
		})
		admin.POST(`type`, createType)

		admin.GET(`permission`, permissions)
	}
	http.ListenAndServe(`:3280`, r)
}
