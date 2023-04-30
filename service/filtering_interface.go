package service

type FilteringServiceInterface interface{
	FilterText ([]string) []string
	checkContains (string) bool
}