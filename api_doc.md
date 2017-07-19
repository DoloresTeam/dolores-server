### **Dolores 后台API 文档**

* 一直打算用文档生成工具来生成这个api文档，不过没有成功。所以暂时就人工先出一个文档吧 

版本：dolores 1.0  
日期：2017-07-19  
作者：Kevin  

API测试地址： http://www.dolores.store:3280

-----

接口列表：

 - 用户登陆获取token

	path:  `/login`  
	method: `post`  
	Parameters: `username` `password`  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	```json
	{
	    "expire": "2017-07-20T11:03:51+08:00",
	    "token": "xxxxx"
	}
	```
	`failure` `httpStatusCode = 401`
	```json
	{
	    "message": "Incorrect Username / Password"
	}
	```


----------

 **在获取Token以后，调用以`api/v1`开头的api都需要在Request Header中添加Authorization 具体格式如下**
 `Authorization: Dolores [token]` 注意Dolores后面有一个空格。用登录接口所获取到的token替换 `[token]`  

 - 刷新用户认证token   
	
	path:  `/api/v1/refresh_token`  
	method: `Get`  
	Parameters: `null`  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{
	    "expire": "2017-07-20T11:03:51+08:00",
	    "token": "xxxxx"
	}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

 - 获取当前登录用户的个人信息

	path:  `/api/v1/profile`  
	method: `Get`  
	Parameters: `null`  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{
	    "cn": "UpDay-github",
	    "email": [],
	    "gender": "1",
	    "id": "b5a62791scgn52v4pqe0",
	    "labeledURI": "FoJFvKRZUpf3fSY9DFnCM0MDCZxM",
	    "name": "Github-UpDay",
	    "priority": "1",
	    "telephoneNumber": "17317309556",
	    "thirdPassword": "v1c+bR]c",
	    "title": "Guest",
	    "unitID": [
	        "b5a60a91scgn4krfr0t0"
	    ]
	}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

- 修改当前登录用户的密码
   
   path:  `/api/v1/modify_password`  
	method: `Post`  
	Parameters: `originalPassword` `newPassword`   原密码 新密码  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

- 通过ID获取一组用户的基本信息

	path:  `/api/v1/basic_profile`  
	method: `Get`  
	Parameters: `id[]`   *用户id数组*  
	Content-Type: `application/json`  
	Result:  
	 `success` `httpStatusCode = 200`
	
	```json
	[
	    {
	        "id": "b5a62791scgn52v4pqe0",
	        "labeledURI": "FoJFvKRZUpf3fSY9DFnCM0MDCZxM",
	        "name": "Github-UpDay"
	    }
	]
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

- 修改当前登录用户的基本信息
	
	path:  `/api/v1/update_profile`  
	method: `Post`  
	Parameters: `labeledURI` `cn` `title`  `email`  头像url、昵称、职位、 邮箱  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

- 获取当前登录用户可见的通讯录视图
	
	path:  `/api/v1/organization`  
	method: `Get`  
	Parameters: `null`  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{
    "departments": [
        {
            "description": "总经办",
            "id": "b4va4h91scghr2vkkaag",
            "ou": "总经办",
            "priority": "1000"
        },
        {
            "description": "事业部",
            "id": "b4va4lh1scghr2vkkab0",
            "ou": "事业部",
            "priority": ""
        },
        {
            "description": "guest from github",
            "id": "b5a60a91scgn4krfr0t0",
            "ou": "Github",
            "priority": "1"
        },
        {
            "description": "人力资源管理",
            "id": "b4va26p1scghr2vkka6g",
            "ou": "人事部",
            "parentID": "b4va1sh1scghr2vkka60",
            "priority": "10"
        }
        ....
    ],
    "members": [
        {
            "cn": "Ew",
            "email": [
                "health.wang@dolores.store"
            ],
            "gender": "0",
            "id": "b4vaqu11scghuujqilgg",
            "labeledURI": "FmIiYdKk0z0CahoJv6M9ZJa6c1qt",
            "name": "王聪灵",
            "priority": "0",
            "telephoneNumber": "18627800585",
            "title": " 项目经理",
            "unitID": [
                "b4va4lh1scghr2vkkab0"
            ]
        },
	      .....
        {
            "cn": "严力",
            "email": [
                "li.yan@dolores.store"
            ],
            "gender": "1",
            "id": "b4vbsnp1scghuujqimhg",
            "labeledURI": "b4vbsnp1scghuujqimhg",
            "name": "严力",
            "priority": "",
            "telephoneNumber": "13000000016",
            "title": "Java工程师",
            "unitID": [
                "b4va7ep1scghr2vkkae0"
            ]
        }
    ],
    "version": "20170719032756Z"
}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`


**注意这里的通讯录视图版本号，客户端需要持久化，可以通过此版本号增量更新通讯录**

----------

- 增量更新本地通讯录视图
	
	path:  `/api/v1/sync_organization／:version`  
	method: `Post`  
	Parameters: `version` *本地通讯录视图版本号*  
	Content-Type: `application/json`  
	Result:
	 `success` `httpStatusCode = 200`
	
	```json
	{
    "logs": [
        {
            "action": "delete",
            "category": "member",
            "content": [
                {
                    "id": "b53pp711scgk1qgv73tg"
                }
            ],
            "createTimestamp": "20170719033843Z"
        },
        {
            "action": "add",
            "category": "member",
            "content": [
                {
                    "cn": "F5",
                    "email": [],
                    "gender": "1",
                    "id": "b5nd8pp1scghhfrajtrg",
                    "labeledURI": "",
                    "name": "F5",
                    "priority": "1",
                    "telephoneNumber": "18888888907",
                    "title": "Guest",
                    "unitID": [
                        "b5a60a91scgn4krfr0t0"
                    ]
                }
            ],
            "createTimestamp": "20170719033919Z"
        },
        {
            "action": "update",
            "category": "member",
            "content": [
                {
                    "cn": "F5",
                    "email": [],
                    "gender": "1",
                    "id": "b5nd8pp1scghhfrajtrg",
                    "labeledURI": "b5nd8pp1scghhfrajtrg",
                    "name": "F5",
                    "priority": "1",
                    "telephoneNumber": "18888888907",
                    "title": "Guest",
                    "unitID": [
                        "b5a60a91scgn4krfr0t0"
                    ]
                }
            ],
            "createTimestamp": "20170719033919Z"
        }
    ],
    "needRefetchOrganization": false,
    "version": "20170719033925Z"
}
	```
	`failure` `httpStatusCode = 401` || `httpStatusCode = 500`

----------

目前后台支持的所有api就这些， 如果大家在使用过程中有任何问题或者建议，欢迎大家加入QQ群讨论交流。群号：`641256202`
