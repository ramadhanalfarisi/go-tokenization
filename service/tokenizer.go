package service

import (
	"regexp"
	"strings"
)

type TokenizerService struct{}

func NewTokenizerService() TokenizerServiceInterface {
	return &TokenizerService{}
}

// CleanString implements TokenizerServiceInterface
func (t *TokenizerService) CreateToken(s string) []string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`(?i)(www\.|https?|s?ftp)\S+`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`[^a-zA-Z0-9-\s]+`).ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	return strings.Fields(s)
}
