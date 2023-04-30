package model

type QueryTokenization struct {
	Text string `form:"text" binding:"required"`
	Language string `form:"lang" binding:"required"`
}
