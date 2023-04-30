package service

type StemmingServiceInterface interface {
	StemText(string) string
	getContainsOrginialWord(string) bool
	removeParticle(string) (string, string)
	removePossesive(string) (string, string)
	removeSuffix(string) (string, string)
	loopSuffix(string, []string) (bool, string)
	removePrefixes(string) (bool, string)
	removePrefix(string) (string, string, []string)
	removePrefixMe(string) (string, []string)
	removePrefixPe(string) (string, []string)
	removePrefixBe(string) (string, []string)
	removePrefixTe(string) (string, []string)
	removeInfix(string) (string, []string)
}
