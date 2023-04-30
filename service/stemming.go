package service

import (
	"strings"

	"github.com/ramadhanalfarisi/go-stopwords-filtering/originalwords"
)

type StemmingService struct {
	Language string
}

func NewStemmingService(lang string) StemmingServiceInterface {
	return &StemmingService{Language: lang}
}

func (s *StemmingService) getContainsOrginialWord(str string) bool {
	if s.Language == "ID" {
		stopword := originalwords.Indonesian
		var _, isexist = stopword[str]
		if isexist {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// StemText implements StemmingServiceInterface
func (s *StemmingService) StemText(str string) string {
	str = strings.ToLower(str)

	var (
		rootFound   = false
		originalstr = str
		particle    string
		possesive   string
		suffix      string
	)

	if len(str) < 3 {
		return str
	}

	if s.getContainsOrginialWord(str) {
		return str
	}

	// Check if prefix must be removed first
	if rxPrefixFirst.MatchString(str) {
		// Remove prefix
		rootFound, str = s.removePrefixes(str)
		if rootFound {
			return str
		}

		// Remove particle
		particle, str = s.removeParticle(str)
		if s.getContainsOrginialWord(str) {
			return str
		}

		// Remove possesive
		possesive, str = s.removePossesive(str)
		if s.getContainsOrginialWord(str) {
			return str
		}

		// Remove suffix
		suffix, str = s.removeSuffix(str)
		if s.getContainsOrginialWord(str) {
			return str
		}
	} else {
		// Remove particle
		particle, str = s.removeParticle(str)
		if s.getContainsOrginialWord(str) {
			return str
		}

		// Remove possesive
		possesive, str = s.removePossesive(str)
		if s.getContainsOrginialWord(str) {
			return str
		}

		// Remove suffix
		suffix, str = s.removeSuffix(str)
		if s.getContainsOrginialWord(str) {
			return str
		}

		// Remove prefix
		rootFound, str = s.removePrefixes(str)
		if rootFound {
			return str
		}

	}

	// If no root found, do loopSuffix
	removedSuffixes := []string{"", suffix, possesive, particle}
	if suffix == "kan" {
		removedSuffixes = []string{"", "k", "an", possesive, particle}
	}

	rootFound, str = s.loopSuffix(originalstr, removedSuffixes)
	if rootFound {
		return str
	}

	// When EVERYTHING failed, return original str
	return originalstr
}

func (s *StemmingService) removeParticle(word string) (string, string) {
	result := rxParticle.ReplaceAllString(word, "")
	particle := strings.Replace(word, result, "", 1)
	return particle, result
}

func (s *StemmingService) removePossesive(word string) (string, string) {
	result := rxPossesive.ReplaceAllString(word, "")
	possesive := strings.Replace(word, result, "", 1)
	return possesive, result
}

func (s *StemmingService) removeSuffix(word string) (string, string) {
	result := rxSuffix.ReplaceAllString(word, "")
	suffix := strings.Replace(word, result, "", 1)
	return suffix, result
}

func (s *StemmingService) loopSuffix(originalWord string, suffixes []string) (bool, string) {
	lenSuffixes := 0
	for _, suffix := range suffixes {
		lenSuffixes += len(suffix)
	}
	wordWithoutSuffix := originalWord[:len(originalWord)-lenSuffixes]

	for i := range suffixes {
		suffixCombination := ""
		for j := 0; j <= i; j++ {
			suffixCombination += suffixes[j]
		}

		word := wordWithoutSuffix + suffixCombination
		if s.getContainsOrginialWord(word) {
			return true, word
		}

		rootFound, word := s.removePrefixes(word)
		if rootFound {
			return true, word
		}
	}

	return false, originalWord
}

func (s *StemmingService) removePrefixes(word string) (bool, string) {
	originalWord := word
	currentPrefix := ""
	removedPrefix := ""
	var recodingChar []string

	for i := 0; i < 3; i++ {
		if len(word) < 3 {
			return false, originalWord
		}

		currentPrefix = word[:2]
		if currentPrefix == removedPrefix {
			break
		}

		removedPrefix, word, recodingChar = s.removePrefix(word)
		if recodingChar != nil {
			for _, char := range recodingChar {
				if s.getContainsOrginialWord(char + word) {
					return true, char + word
				}
			}

			if s.getContainsOrginialWord(word) {
				return true, word
			}
		} else {
			if s.getContainsOrginialWord(word) {
				return true, word
			}
		}

	}

	return false, word
}

func (s *StemmingService) removePrefix(word string) (prefix string, result string, recoding []string) {
	if strings.HasPrefix(word, "kau") {
		return "kau", word[3:], nil
	}

	prefix = word[:2]
	switch prefix {
	case "di", "ke", "se", "ku":
		result = word[2:]
	case "me":
		result, recoding = s.removePrefixMe(word)
	case "pe":
		result, recoding = s.removePrefixPe(word)
	case "be":
		result, recoding = s.removePrefixBe(word)
	case "te":
		result, recoding = s.removePrefixTe(word)
	default:
		result, recoding = s.removeInfix(word)
	}

	return prefix, result, recoding
}

func (s *StemmingService) removePrefixMe(word string) (string, []string) {
	// Pattern 1
	// me{l|r|w|y}V => me-{l|r|w|y}V
	matches := rxPrefixMe1.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 2
	// mem{b|f|v} => mem-{b|f|v}
	matches = rxPrefixMe2.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 3
	// mempe => mem-pe
	matches = rxPrefixMe3.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 4
	// mem{rV|V} => mem-{rV|V} OR me-p{rV|V}
	matches = rxPrefixMe4.FindStringSubmatch(word)
	if len(matches) != 0 {
		if matches[1] == "esan" || matches[1] == "uas" {
			return matches[1], []string{"p"}
		} else {
			return matches[1], []string{"m", "p"}
		}

	}

	// Pattern 5
	// men{c|d|j|s|t|z} => men-{c|d|j|s|t|z}
	matches = rxPrefixMe5.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 6
	// menV => nV OR tV
	matches = rxPrefixMe6.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"n", "t"}
	}

	// Pattern 7
	// meng{g|h|q|k} => meng-{g|h|q|k}
	matches = rxPrefixMe7.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 8
	// mengV => meng-V OR meng-kV OR me-ngV OR mengV- where V = 'e'
	matches = rxPrefixMe8.FindStringSubmatch(word)
	if len(matches) != 0 {
		if matches[2] == "e" {
			return matches[3], nil
		}

		return matches[1], []string{"ng", "k"}
	}

	// Pattern 9
	// menyV => meny-sV OR me-nyV to stem menyala
	matches = rxPrefixMe9.FindStringSubmatch(word)
	if len(matches) != 0 {
		if matches[2] == "a" {
			return "ny" + matches[1], nil
		}

		return "s" + matches[1], nil
	}

	// Pattern 10
	// mempV => mem-pA where A != 'e'
	matches = rxPrefixMe10.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	return word, nil
}

func (s *StemmingService) removePrefixPe(word string) (string, []string) {
	// Pattern 1
	// pe{w|y}V => pe-{w|y}V
	matches := rxPrefixPe1.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 2
	// perV => per-V OR pe-rV
	matches = rxPrefixPe2.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"r"}
	}

	// Pattern 3
	// perCAP => per-CAP where C != 'r' and P != 'er'
	matches = rxPrefixPe3.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 4
	// perCAerV => per-CAerV where C != 'r'
	matches = rxPrefixPe4.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 5
	// pem{b|f|v} => pem-{b|f|v}
	matches = rxPrefixPe5.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 6
	// pem{rV|V} => pe-m{rV|V} OR pe-p{rV|V}
	matches = rxPrefixPe6.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"m", "p"}
	}

	// Pattern 7
	// pen{c|d|j|s|t|z} => pen-{c|d|j|s|t|z}
	matches = rxPrefixPe7.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 8
	// penV => pe-nV OR pe-tV
	matches = rxPrefixPe8.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"n", "t"}
	}

	// Pattern 9
	// pengC => peng-C
	matches = rxPrefixPe9.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 10
	// pengV => peng-V OR peng-kV OR pengV- where V = 'e'
	matches = rxPrefixPe10.FindStringSubmatch(word)
	if len(matches) != 0 {
		if matches[2] == "e" {
			return matches[3], nil
		}

		return matches[1], []string{"k"}
	}

	// Pattern 11
	// penyV => peny-sV OR pe-nyV
	matches = rxPrefixPe11.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"s", "ny"}
	}

	// Pattern 12
	// pelV => pe-lV OR pel-V for pelajar
	matches = rxPrefixPe12.FindStringSubmatch(word)
	if len(matches) != 0 {
		if word == "pelajar" {
			return "ajar", nil
		}

		return matches[1], nil
	}

	// Pattern 13
	// peCerV => peC-erV where C != {r|w|y|l|m|n}
	matches = rxPrefixPe13.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 14
	// peCP => pe-CP where C != {r|w|y|l|m|n} and P != 'er'
	matches = rxPrefixPe14.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 15
	// peC1erC2 => pe-C1erC2 where C1 != {r|w|y|l|m|n}
	matches = rxPrefixPe15.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	return word, nil
}

func (s *StemmingService) removePrefixBe(word string) (string, []string) {
	// Pattern 1
	// berV => ber-V OR be-rV
	matches := rxPrefixBe1.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"r"}
	}

	// Pattern 2
	// berCAP => ber-CAP where C != 'r' and P != 'er'
	matches = rxPrefixBe2.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 3
	// berCAerV => ber-CAerV where C != 'r'
	matches = rxPrefixBe3.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 4
	// belajar => bel-ajar
	matches = rxPrefixBe4.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 5
	// beC1erC2 => be-C1erC2 where C1 != {'r'|'l'}
	matches = rxPrefixBe5.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	return word, nil
}

func (s *StemmingService) removePrefixTe(word string) (string, []string) {
	// Pattern 1
	// terV => ter-V OR te-rV
	matches := rxPrefixTe1.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], []string{"r"}
	}

	// Pattern 2
	// terCerV => ter-CerV where C != 'r'
	matches = rxPrefixTe2.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 3
	// terCP => ter-CP where C != 'r' and P != 'er'
	matches = rxPrefixTe3.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 4
	// teC1erC2 => te-C1erC2 where C1 != 'r'
	matches = rxPrefixTe4.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	// Pattern 5
	// terC1erC2 => ter-C1erC2 where C1 != 'r'
	matches = rxPrefixTe5.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[1], nil
	}

	return word, nil
}

func (s *StemmingService) removeInfix(word string) (string, []string) {
	// Pattern 1
	// CerV => CerV OR CV
	matches := rxInfix1.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[3], []string{matches[1], matches[2]}
	}

	// Pattern 2
	// CinV => CinV OR CV
	matches = rxInfix2.FindStringSubmatch(word)
	if len(matches) != 0 {
		return matches[3], []string{matches[1], matches[2]}
	}

	return word, nil
}
