# dolores-server
dolores 服务端 提供RestFull API

## 接口列表

### /organization
###### 参数: 无
###### 返回值:
``` json
{
  "departments": [ // 返回当前登录用户所有可见的部门， 父部门总是早于子部门返回
    {
      "description": "This is test unit's description",
      "id": "b4a1bldhfpcl3m5lhohg",
      "name": "Test",
    },
    {
      "description": "This is test unit's description", // 部门的描述备注信息
      "id": "b49kdrg6h302hrpggg8g", // 部门ID
      "name": "Test", // 部门名称
      "pid": "b4a1bldhfpcl3m5lhohg" // 父部门ID
    },
    {
      "description": "This is test unit's description",
      "id": "b4a1cqlhfpcl6964kqtg",
      "name": "Test",
      "pid": "b49kdrg6h302hrpggg8g"
    }
  ],
  "members": [ // 返回当前登录用户的所有可见员工
    {
      "departmentIDs": [ // 员工所在部门id 一个员工可能在多个部门
        "b49kdrg6h302hrpggg8g"
      ],
      "easemobAccount": "", // 环信ID
      "easemobPassword": "", // 环信密码
      "email": [ // 用户邮箱，可能有多个邮箱
        "aoxianglele@icloud.com"
      ],
      "id": "b49kehg6h302jg98oi70", // 用户ID
      "name": "Kevin.Gong", // 显示用的名字
      "realName": "巩祥", // 真名
      "title": [ // 职位，可能有多个
        "Developer"
      ],
      "avatarURL": "http://www.dolores.store/static/avatar/b49kehg6h302jg98oi70.png", // 用户头像url
      "gender": 0, // 用户性别 0: 女 1: 男 其他: 未知
    }
  ],
  "version": 1
}
```
