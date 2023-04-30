package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ramadhanalfarisi/go-stopwords-filtering/model"
	"github.com/ramadhanalfarisi/go-stopwords-filtering/service"
)

type TokenizationController struct{}

func NewTokenizationController() TokenizationControllerInterface {
	return &TokenizationController{}
}

// StemText implements TokenizationControllerInterface
func (s *TokenizationController) StemText(g *gin.Context) {
	var query model.QueryTokenization

	if g.ShouldBind(&query) == nil {
		var result []string
		output := make(chan string)
		sservice := service.NewFilteringService(query.Language)
		tservice := service.NewTokenizerService()
		stservice := service.NewStemmingService(query.Language)
		as := tservice.CreateToken(query.Text)
		res := sservice.FilterText(as)
		for _, item := range res {
			go func(str string) {
				str = stservice.StemText(str)
				output <- str
			}(item)
		}
		for i := 0; i < len(res); i++ {
			str := <-output
			result = append(result, str)
		}
		g.JSON(200, gin.H{"status": "success", "data": result})
	} else {
		g.JSON(500, gin.H{"status": "failed", "data": nil})
	}
}

// TokenizeText implements TokenizationControllerInterface
func (s *TokenizationController) TokenizeText(g *gin.Context) {
	var query model.QueryTokenization

	if g.ShouldBind(&query) == nil {
		tservice := service.NewTokenizerService()
		as := tservice.CreateToken(query.Text)
		g.JSON(200, gin.H{"status": "success", "data": as})
	} else {
		g.JSON(500, gin.H{"status": "failed", "data": nil})
	}
}

// FilterText implements TokenizationControllerInterface
func (*TokenizationController) FilterText(g *gin.Context) {
	var query model.QueryTokenization

	if g.ShouldBind(&query) == nil {
		sservice := service.NewFilteringService(query.Language)
		tservice := service.NewTokenizerService()
		as := tservice.CreateToken(query.Text)
		res := sservice.FilterText(as)
		g.JSON(200, gin.H{"status": "success", "data": res})
	} else {
		g.JSON(500, gin.H{"status": "failed", "data": nil})
	}
}
