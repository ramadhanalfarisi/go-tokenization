package main

import 	"github.com/ramadhanalfarisi/go-stopwords-filtering/app"


func main(){
	var a app.App
	a.CreateRoutes()
	a.Run()
}