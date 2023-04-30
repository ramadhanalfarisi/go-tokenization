package service

import (

	"github.com/ramadhanalfarisi/go-stopwords-filtering/stopwords"
)

type FilteringService struct{
	Language string
}

func NewFilteringService(lang string) FilteringServiceInterface {
	return &FilteringService{Language: lang}
}

// CheckContains implements FilteringServiceInterface
func (s *FilteringService) checkContains(text string) bool {
	if s.Language == "ID" {
		stopword := stopwords.Indonesian
		var _, isexist = stopword[text]
		if isexist {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// FilterText implements FilteringServiceInterface
func (s *FilteringService) FilterText(as []string) []string {
	var result []string
	for _, item := range as {
		if iscontains := s.checkContains(item); !iscontains {
			result = append(result, item)
		}
	}
	return result
}

