# dolores-server

dolores 服务端 提供REST API

## API 测试地址
```
  http://www.dolores.store:3280
```

## 接口列表

#### **客户端**

登录：  `[POST]`  ` /login`  
参数： body中用json编码 username password  
备注：下面所有 `/api/v1` 路径下的接口都需要添加授权头，具体做法：
在 `http header` 中添加 `Authorization=Dolores + 登录接口返回的Token`(注意Dolores后面有一个空格) 

----------
刷新Token： `[GET]`   `/api/v1/resfresh_token`  
参数：无  
备注: `Token` 的有效期为24小时，最大刷新间隔为1星期  

----------
获取个人信息： `[GET]`  `/api/v1/profile`  
参数：无  
备注:  

----------
更新头像：`[GET]` `/api/v1/update_avatar`  
参数：body中用json编码 avatar  
备注: 需要客户端自己上传图片到七牛，然后将上传后的`url`回传到服务器

----------
获取组织架构：`[GET]` `/api/v1/organization`  
参数：无  
备注:  

#### **管理后台**
