package main

import (
	"github.com/gin-gonic/gin"
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
