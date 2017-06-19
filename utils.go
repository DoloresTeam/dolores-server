package main

import (
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
)

// PutRet  构造返回值字段
type PutRet struct {
	Hash     string `json:"hash"`
	Key      string `json:"key"`
	Filesize int    `json:"filesize"`
}

func uploadFileToQiNiu(filePath, key string) (*PutRet, error) {
	// 创建一个Client
	c := kodo.New(0, nil)
	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: `dolores:` + key,
		//设置Token过期时间
		Expires: 3600,
	}
	// 生成一个上传token
	token, err := c.MakeUptokenWithSafe(policy)
	if err != nil {
		return nil, err
	}
	// 构建一个uploader
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)
	var ret PutRet
	// 调用PutFile方式上传，这里的key需要和上传指定的key一致
	err = uploader.PutFile(nil, &ret, token, key, filePath, nil)
	return &ret, err
}
