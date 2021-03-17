package routers

import "github.com/gin-gonic/gin"

type Option func(e *gin.Engine)

var options []Option

func Init() *gin.Engine {
	r := gin.Default()
	for _, opt := range options {
		opt(r)
	}
	return r
}

func Include(opts ...Option) {
	options = append(options, opts...)
}
