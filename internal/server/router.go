package server

import "github.com/gin-gonic/gin"

// Router interface for service route registration
type Router interface {
	Register(engine gin.IRouter)
}
