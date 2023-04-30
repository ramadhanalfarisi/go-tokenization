package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ramadhanalfarisi/go-stopwords-filtering/controller"
)

type App struct {
	Routes *gin.Engine
}

func (a *App) CreateRoutes(){
	g := gin.Default()
	controller := controller.NewTokenizationController()
	g.GET("/filtering", controller.FilterText)
	g.GET("/tokenize", controller.TokenizeText)
	g.GET("/stemming", controller.StemText)
	a.Routes = g
}

func (a *App) Run(){
	a.Routes.Run(":8080")
}

