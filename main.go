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
		admin.GET(`type`, fetchTypes)
		admin.GET(`type/:id`, typeByID)
		admin.PUT(`type/:id`, editType)
		admin.DELETE(`type/:id`, delType)
		admin.POST(`type`, createType)

		admin.GET(`permission`, fetchPermissions)
		admin.GET(`u_permission`, fetchPermissionByRoleUPID)
		admin.GET(`p_permission`, fetchPermissionByRolePPID)
		admin.GET(`permission/:id`, permissionByID)
		admin.PUT(`permission/:id`, editPermission)
		admin.DELETE(`permission/:id`, delPermission)
		admin.POST(`permission`, createPermission)

		admin.GET(`role`, fetchRoles)
		admin.POST(`role`, createRole)
		admin.GET(`role/:id`, roleByID)
		admin.PUT(`role/:id`, editRole)
		admin.DELETE(`role/:id`, delRole)

		admin.GET(`member`, fetchMembers)

		admin.GET(`department`, fetchDepartment)
		admin.POST(`department`, createDepartment)
		admin.GET(`department/:id`, departmentByID)
	}

	http.ListenAndServe(`:3280`, r)

}
