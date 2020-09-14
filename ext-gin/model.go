package extgin

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/karsto/glew/common/ext-gin/controllers"
	"github.com/karsto/glew/common/ext-gin/middleware"
)

type Registerer interface {
	Register(s *gin.RouterGroup)
}

type ControllerGroup interface {
	GetControllers() []Registerer
}

func NewStd() *gin.Engine {
	s := gin.New()

	s.Use(gin.Recovery()) // should likely be first .Use(middleware)
	s.Use(gin.ErrorLogger())
	s.Use(middleware.ErrorBodyLogger())
	s.Use(middleware.ResponseLogger(false)) // TODO: env
	s.Use(gin.Logger())
	s.Use(gzip.Gzip(gzip.DefaultCompression))
	// s.Use(middleware.RequestLogger())  // TODO:  env

	s.RedirectTrailingSlash = true // removes gin's weird trailing slash sensitivity related, https://github.com/gin-gonic/gin/pull/1061
	s.RedirectFixedPath = true     // makes the best damn attempt it can to match paths (case-insensitive and trimmed)

	statusController := controllers.NewStatusController()
	statusController.Register(s.Group(""))

	return s
}
