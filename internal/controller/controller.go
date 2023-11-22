package controller

import (
	"net/http"

	"example.com/dynamicWordpressBuilding/internal/service"
	"example.com/dynamicWordpressBuilding/utils"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Router     *gin.Engine
	svc        service.ServiceInterface
	tokenMaker utils.Maker
}

func NewController(svc service.ServiceInterface) *Controller {
	ctl := &Controller{}
	ctl.Router = gin.Default()
	ctl.svc = svc
	ctl.tokenMaker = utils.NewTokenMaker()
	ctl.Routes()
	return ctl
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, "this is a page for dynamic website building")
}
