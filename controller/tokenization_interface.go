package controller

import "github.com/gin-gonic/gin"

type TokenizationControllerInterface interface {
	TokenizeText(*gin.Context)
	FilterText(*gin.Context)
	StemText(*gin.Context)
}