package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusController struct {
}

func NewStatusController() *StatusController {
	c := StatusController{}
	return &c
}

func (c *StatusController) Register(service *gin.RouterGroup) {
	router := service.Group("status")
	{
		router.GET("/", c.Status)
	}
}

// TODO: app info things like uptime
// TODO: add metrics endpoint https://github.com/gin-contrib/expvar
// TODO: expose function to add StatusCheck and loop through them
// TODO: common status checks: db connectivity, db exists, tables exist, internet, twilio, aws resources.
// TODO: errors endpoint returning 500 / commons for stimulating / testing  - should return example error body structures with metadata
// TODO: status websocket for connectivity checks
type StatusCheck func() (name string, msg string, err error)

// Status - an unauthenticated endpoint to verify the api is running
// @Summary returns a status message
// @Description returns a status message
// @Tags status
// @Accept  json
// @Produce  json
// @Success 200 {string} string ""
// @Router /status/ [get]
func (*StatusController) Status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "up up up",
	},
	)
}
