package service

import (
	"unicode"
)

type StemmingEnService struct{}

func NewStemmingEnService() StemmingEnServiceInterface {
	return &StemmingEnService{}
}

func (stem *StemmingEnService) isConsonant(s []rune, i int) bool {

	//DEBUG
	//log.Printf("isConsonant: [%+v]", string(s[i]))

	result := true

	switch s[i] {
	case 'a', 'e', 'i', 'o', 'u':
		result = false
	case 'y':
		if i == 0 {
			result = true
		} else {
			result = !stem.isConsonant(s, i-1)
		}
	default:
		result = true
	}

	return result
}

func (stem *StemmingEnService) measure(s []rune) uint {

	// Initialize.
	lenS := len(s)
	result := uint(0)
	i := 0

	// Short Circuit.
	if lenS == 0 {
		/////////// RETURN
		return result
	}

	// Ignore (potential) consonant sequence at the beginning of word.
	for stem.isConsonant(s, i) {

		//DEBUG
		//log.Printf("[measure([%s])] Eat Consonant [%d] -> [%s]", string(s), i, string(s[i]))

		i++
		if i >= lenS {
			/////////////// RETURN
			return result
		}
	}

	// For each pair of a vowel sequence followed by a consonant sequence, increment result.
	Outer:
	for i < lenS {

		for !stem.isConsonant(s, i) {

			//DEBUG
			//log.Printf("[measure([%s])] VOWEL [%d] -> [%s]", string(s), i, string(s[i]))

			i++
			if i >= lenS {
				/////////// BREAK
				break Outer
			}
		}
		for stem.isConsonant(s, i) {

			//DEBUG
			//log.Printf("[measure([%s])] CONSONANT [%d] -> [%s]", string(s), i, string(s[i]))

			i++
			if i >= lenS {
				result++
				/////////// BREAK
				break Outer
			}
		}
		result++
	}

	// Return
	return result
}

func (stem *StemmingEnService) hasSuffix(s, suffix []rune) bool {

	lenSMinusOne := len(s) - 1
	lenSuffixMinusOne := len(suffix) - 1

	if lenSMinusOne <= lenSuffixMinusOne {
		return false
	} else if s[lenSMinusOne] != suffix[lenSuffixMinusOne] { // I suspect checking this first should speed this func (stem *StemmingEnService)tion up in practice.
		/////// RETURN
		return false
	} else {

		for i := 0; i < lenSuffixMinusOne; i++ {

			if suffix[i] != s[lenSMinusOne-lenSuffixMinusOne+i] {
				/////////////// RETURN
				return false
			}

		}

	}

	return true
}

func (stem *StemmingEnService) containsVowel(s []rune) bool {

	lenS := len(s)

	for i := 0; i < lenS; i++ {

		if !stem.isConsonant(s, i) {
			/////////// RETURN
			return true
		}

	}

	return false
}

func (stem *StemmingEnService) hasRepeatDoubleConsonantSuffix(s []rune) bool {

	// Initialize.
	lenS := len(s)

	result := false

	// Do it!
	if lenS < 2 {
		result = false
	} else if s[lenS-1] == s[lenS-2] && stem.isConsonant(s, lenS-1) { // Will using isConsonant() cause a problem with "YY"?
		result = true
	} else {
		result = false
	}

	// Return,
	return result
}

func (stem *StemmingEnService) hasConsonantVowelConsonantSuffix(s []rune) bool {

	// Initialize.
	lenS := len(s)

	result := false

	// Do it!
	if lenS < 3 {
		result = false
	} else if stem.isConsonant(s, lenS-3) && !stem.isConsonant(s, lenS-2) && stem.isConsonant(s, lenS-1) {
		result = true
	} else {
		result = false
	}

	// Return
	return result
}

func (stem *StemmingEnService) step1a(s []rune) []rune {

	// Initialize.
	var result []rune = s

	lenS := len(s)

	// Do it!
	if suffix := []rune("sses"); stem.hasSuffix(s, suffix) {

		lenTrim := 2

		subSlice := s[:lenS-lenTrim]

		result = subSlice
	} else if suffix := []rune("ies"); stem.hasSuffix(s, suffix) {
		lenTrim := 2

		subSlice := s[:lenS-lenTrim]

		result = subSlice
	} else if suffix := []rune("ss"); stem.hasSuffix(s, suffix) {

		result = s
	} else if suffix := []rune("s"); stem.hasSuffix(s, suffix) {

		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]

		result = subSlice
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step1b(s []rune) []rune {

	// Initialize.
	var result []rune = s

	lenS := len(s)

	// Do it!
	if suffix := []rune("eed"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 0 {
			lenTrim := 1

			result = s[:lenS-lenTrim]
		}
	} else if suffix := []rune("ed"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		if stem.containsVowel(subSlice) {

			if suffix2 := []rune("at"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if suffix2 := []rune("bl"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if suffix2 := []rune("iz"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; c != 'l' && c != 's' && c != 'z' && stem.hasRepeatDoubleConsonantSuffix(subSlice) {
				lenTrim := 1

				lenSubSlice := len(subSlice)

				result = subSlice[:lenSubSlice-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; stem.measure(subSlice) == 1 && stem.hasConsonantVowelConsonantSuffix(subSlice) && c != 'w' && c != 'x' && c != 'y' {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else {
				result = subSlice
			}

		}
	} else if suffix := []rune("ing"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		if stem.containsVowel(subSlice) {

			if suffix2 := []rune("at"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if suffix2 := []rune("bl"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if suffix2 := []rune("iz"); stem.hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if c := subSlice[len(subSlice)-1]; c != 'l' && c != 's' && c != 'z' && stem.hasRepeatDoubleConsonantSuffix(subSlice) {
				lenTrim := 1

				lenSubSlice := len(subSlice)

				result = subSlice[:lenSubSlice-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; stem.measure(subSlice) == 1 && stem.hasConsonantVowelConsonantSuffix(subSlice) && c != 'w' && c != 'x' && c != 'y' {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else {
				result = subSlice
			}

		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step1c(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	result := s

	// Do it!
	if 2 > lenS {
		/////////// RETURN
		return result
	}

	if s[lenS-1] == 'y' && stem.containsVowel(s[:lenS-1]) {

		result[lenS-1] = 'i'

	} else if s[lenS-1] == 'Y' && stem.containsVowel(s[:lenS-1]) {

		result[lenS-1] = 'I'

	}

	// Return.
	return result
}

func (stem *StemmingEnService) step2(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	result := s

	// Do it!
	if suffix := []rune("ational"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-5] = 'e'
			result = result[:lenS-4]
		}
	} else if suffix := []rune("tional"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = result[:lenS-2]
		}
	} else if suffix := []rune("enci"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("anci"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("izer"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-1]
		}
	} else if suffix := []rune("bli"); stem.hasSuffix(s, suffix) { // --DEPARTURE--
		//		} else if suffix := []rune("abli") ; stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("alli"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("entli"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("eli"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("ousli"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("ization"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-5] = 'e'

			result = s[:lenS-4]
		}
	} else if suffix := []rune("ation"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-3] = 'e'

			result = s[:lenS-2]
		}
	} else if suffix := []rune("ator"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-2] = 'e'

			result = s[:lenS-1]
		}
	} else if suffix := []rune("alism"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-3]
		}
	} else if suffix := []rune("iveness"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("fulness"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("ousness"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("aliti"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result = s[:lenS-3]
		}
	} else if suffix := []rune("iviti"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-3] = 'e'

			result = result[:lenS-2]
		}
	} else if suffix := []rune("biliti"); stem.hasSuffix(s, suffix) {
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			result[lenS-5] = 'l'
			result[lenS-4] = 'e'

			result = result[:lenS-3]
		}
	} else if suffix := []rune("logi"); stem.hasSuffix(s, suffix) { // --DEPARTURE--
		if stem.measure(s[:lenS-len(suffix)]) > 0 {
			lenTrim := 1

			result = s[:lenS-lenTrim]
		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step3(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if suffix := []rune("icate"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if stem.measure(s[:lenS-lenSuffix]) > 0 {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("ative"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 0 {
			result = subSlice
		}
	} else if suffix := []rune("alize"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if stem.measure(s[:lenS-lenSuffix]) > 0 {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("iciti"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if stem.measure(s[:lenS-lenSuffix]) > 0 {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("ical"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if stem.measure(s[:lenS-lenSuffix]) > 0 {
			result = result[:lenS-2]
		}
	} else if suffix := []rune("ful"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 0 {
			result = subSlice
		}
	} else if suffix := []rune("ness"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 0 {
			result = subSlice
		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step4(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if suffix := []rune("al"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("ance"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("ence"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("er"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ic"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("able"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ible"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ant"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ement"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ment"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ent"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ion"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		c := subSlice[len(subSlice)-1]

		if m > 1 && (c == 's' || c == 't') {
			result = subSlice
		}
	} else if suffix := []rune("ou"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ism"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ate"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("iti"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ous"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ive"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	} else if suffix := []rune("ize"); stem.hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step5a(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if s[lenS-1] == 'e' {
		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]
		if len(subSlice) == 0 {
			return result
		}
		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		} else if m == 1 {
			if c := subSlice[len(subSlice)-1]; !(stem.hasConsonantVowelConsonantSuffix(subSlice) && c != 'w' && c != 'x' && c != 'y') {
				result = subSlice
			}
		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) step5b(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if 2 < lenS && s[lenS-2] == 'l' && s[lenS-1] == 'l' {

		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]

		m := stem.measure(subSlice)

		if m > 1 {
			result = subSlice
		}
	}

	// Return.
	return result
}

func (stem *StemmingEnService) StemEnText(str string) string {

	// Convert string to []rune
	runeArr := []rune(str)

	// Stem.
	runeArr = stem.Stem(runeArr)

	// Convert []rune to string
	strg := string(runeArr)

	// Return.
	return strg
}

func (stem *StemmingEnService) Stem(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	// Short circuit.
	if lenS == 0 {
		/////////// RETURN
		return s
	}

	// Make all runes lowercase.
	for i := 0; i < lenS; i++ {
		s[i] = unicode.ToLower(s[i])
	}

	// Stem
	result := stem.StemWithoutLowerCasing(s)

	// Return.
	return result
}

func (stem *StemmingEnService) StemWithoutLowerCasing(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	// Words that are of length 2 or less is already stemmed.
	// Don't do anything.
	if lenS <= 2 {
		/////////// RETURN
		return s
	}

	// Stem
	s = stem.step1a(s)
	s = stem.step1b(s)
	s = stem.step1c(s)
	s = stem.step2(s)
	s = stem.step3(s)
	s = stem.step4(s)
	s = stem.step5a(s)
	s = stem.step5b(s)

	// Return.
	return s
}
