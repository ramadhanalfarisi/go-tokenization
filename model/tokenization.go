package model

type QueryTokenization struct {
	Text string `form:"text"`
	Language string `form:"lang"`
}
