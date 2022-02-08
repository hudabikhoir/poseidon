package api

import (
	"boilerplate-golang-v2/api/v1/content"

	"github.com/labstack/echo"
)

// Controller to define controller that we use
type Controller struct {
	ContentController *content.Controller
}

//RegisterPath Registera V1 API path
func RegisterPath(e *echo.Echo, ctrl Controller) {
	//content
	contentV1 := e.Group("v1/contents")
	contentV1.GET("/", ctrl.ContentController.FindContentByTag)
	contentV1.GET("/:id", ctrl.ContentController.GetContentByID)
	contentV1.GET("/tag/:tag", ctrl.ContentController.FindContentByTag)
	contentV1.POST("", ctrl.ContentController.CreateNewContent)
	contentV1.PUT("/:id", ctrl.ContentController.UpdateContent)

	//health check
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})
}
