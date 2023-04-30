package service

type TokenizerServiceInterface interface {
	CreateToken(string) []string
}
