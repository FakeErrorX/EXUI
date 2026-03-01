package controller

import (
	"github.com/gin-gonic/gin"
)

type EXUIController struct {
	BaseController

	inboundController     *InboundController
	settingController     *SettingController
	xraySettingController *XraySettingController
}

func NewEXUIController(g *gin.RouterGroup) *EXUIController {
	a := &EXUIController{}
	a.initRouter(g)
	return a
}

func (a *EXUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/panel")
	g.Use(a.checkLogin)

	g.GET("/", a.index)
	g.GET("/inbounds", a.inbounds)
	g.GET("/settings", a.settings)
	g.GET("/xray", a.xraySettings)

	a.inboundController = NewInboundController(g)
	a.settingController = NewSettingController(g)
	a.xraySettingController = NewXraySettingController(g)
}

func (a *EXUIController) index(c *gin.Context) {
	html(c, "index.html", "pages.index.title", nil)
}

func (a *EXUIController) inbounds(c *gin.Context) {
	html(c, "inbounds.html", "pages.inbounds.title", nil)
}

func (a *EXUIController) settings(c *gin.Context) {
	html(c, "settings.html", "pages.settings.title", nil)
}

func (a *EXUIController) xraySettings(c *gin.Context) {
	html(c, "xray.html", "pages.xray.title", nil)
}
